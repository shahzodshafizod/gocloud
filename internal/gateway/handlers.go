package gateway

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"os"
	"strings"

	"github.com/pkg/errors"
	_ "github.com/shahzodshafizod/gocloud/docs"
	"github.com/shahzodshafizod/gocloud/internal/response"
	"github.com/shahzodshafizod/gocloud/internal/validator"
	"github.com/shahzodshafizod/gocloud/pkg"
	"github.com/shahzodshafizod/gocloud/pkg/http"
	httpSwagger "github.com/swaggo/http-swagger"
	"go.uber.org/fx"
)

const (
	_USER_ID_KEY      = "userID"
	_MAX_AVATAR_BYTES = 5 << 20 // 101 << 20 => 10100000000000000000000 => 2^22+2^20 => 5242880/1024/1024 = 5mb
)

type handler struct {
	service      Service
	apiSecretKey string
}

func NewHandler(
	lifecycle fx.Lifecycle,
	service Service,
	tracer pkg.Tracer,
	queue pkg.Queue,
) {
	handler := &handler{
		service:      service,
		apiSecretKey: os.Getenv("API_SECRET_KEY"),
	}

	router := handler.registerRoutes(tracer)

	lifecycle.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			go router.Serve(ctx)
			return nil
		},
		OnStop: func(ctx context.Context) error {
			tracer.Shutdown(ctx)
			queue.Close()
			return router.Shutdown(ctx)
		},
	})
}

func (h *handler) registerRoutes(tracer pkg.Tracer) pkg.Router {
	var (
		addr      = os.Getenv("SERVICE_ADDRESS")
		prefixFmt = "/api/v1/%s" // %s - group
	)

	router := http.New(addr, tracer, validator.OptionalDateOnly())

	router.GET("/docs/", "Swagger", router.WrapHandler(httpSwagger.WrapHandler))

	prefix := fmt.Sprintf(prefixFmt, "users")
	router.POST(prefix+"/signup", "SignUp", h.signUp)
	router.POST(prefix+"/signup/confirm", "ConfirmSignUp", h.confirmSignUp)
	router.POST(prefix+"/signin", "SignIn", h.signIn)
	router.POST(prefix+"/token/refresh/:userid", "RefreshToken", h.refreshToken) // TODO: secure this route
	router.GET(prefix+"/profile", "GetUser", h.getUser, h.authorize)
	router.PUT(prefix+"/update", "UpdateUser", h.updateUser, h.authorize)
	router.PUT(prefix+"/update/confirm", "ConfirmChangeEmail", h.confirmChangeEmail)
	router.PUT(prefix+"/update/avatar", "UpdateAvatar", h.updateAvatar, h.authorize)
	router.GET(prefix+"/password/forgot", "ForgotPassword", h.forgotPassword, h.authorize)
	router.PUT(prefix+"/password/reset", "ResetPassword", h.resetPassword)
	router.PUT(prefix+"/password/change", "ChangePassword", h.changePassword, h.authorize)
	router.POST(prefix+"/signout", "SignOut", h.signOut, h.authorize)
	router.DELETE(prefix+"/delete", "DeleteUser", h.deleteUser, h.authorize)

	prefix = fmt.Sprintf(prefixFmt, "partners")
	router.GET(prefix+"/products", "GetProducts", h.getPartnerProducts, h.authorize)

	prefix = fmt.Sprintf(prefixFmt, "orders")
	router.POST(prefix+"/check", "CheckOrder", h.checkOrder, h.authorize)
	router.POST(prefix+"/confirm", "ConfirmOrder", h.confirmOrder, h.authorize)
	router.POST(prefix+"/pay", "PayOrder", h.payOrder)
	router.POST(prefix+"/pickup", "PickUpOrder", h.pickUpOrder)
	router.POST(prefix+"/assign", "AssignOrder", h.assignOrder, h.authorize)

	return router
}

// SignUp godoc
//
//	@Summary		User SignUp
//	@Tags			users-auth
//	@Description	create user account
//	@Accept			json
//	@Produce		json
//	@Param			Signature	header		string	true	"hash_hmac('sha256', first_name+last_name+email+password+role, apiSecretKey)"
//	@Param			Request		body		signUp	true	"sign up info"
//	@Success		200			{object}	response.response{payload=signUpResponse}
//	@Failure		400			"Bad Request"
//	@Failure		401			"Unauthorized"
//	@Failure		404			"Not Found"
//	@Failure		409			"Already Exists"
//	@Failure		410			"Wrong Verification"
//	@Failure		411			"Wrong Password"
//	@Failure		412			"Invalid Token"
//	@Failure		415			"Unsupported Avatar Format"
//	@Failure		500			"Internal Server Error"
//	@Failure		504			"External Service Error"
//	@Router			/users/signup [post]
//	@Security		Request Signature
func (h *handler) signUp(c pkg.Context) {
	ctx, span := c.StartSpan()
	defer span.End()

	var signUp = &signUp{}
	err := c.ParseBody(signUp)
	if err != nil {
		c.Respond(response.Make(response.BadRequestCode).WithMessage("ParseBody: " + err.Error()))
		return
	}

	errs := c.ValidateStruct(signUp)
	if len(errs) != 0 {
		c.Respond(response.Make(response.BadRequestCode).WithMessage("ValidateStruct: " + strings.Join(errs, "; ")))
		return
	}

	signature, err := generateHashByKey(h.apiSecretKey,
		signUp.FirstName+signUp.LastName+signUp.Email+signUp.Password+signUp.Role)
	if err != nil {
		span.RecordError(errors.Wrap(err, "generateHashByKey"))
		c.Respond(response.Make(response.InternalServerErrorCode).WithMessage(err.Error()))
		return
	}

	if c.GetHeader("Signature") != signature {
		span.RecordError(errors.New("wrong Signature: " + signature))
		c.Respond(response.Make(response.UnauthorizedCode).WithMessage("wrong signup request"))
		return
	}

	resp, err := h.service.signUp(ctx, signUp)
	if err != nil {
		span.RecordError(errors.Wrap(err, "h.service.signUp"))
		c.Respond(response.Make(response.BadGatewayCode))
		return
	}

	c.Respond(response.Make(response.OKCode).WithPayload(resp))
}

// ConfirmSignUp godoc
//
//	@Summary		Confirm User SignUp
//	@Tags			users-auth
//	@Description	Verifies a user's email address using a verification code. Activates the user's account after successful verification.
//	@Accept			json
//	@Produce		json
//	@Param			Request	body		confirmSignUp	true	"confirm sign up info"
//	@Success		200		{object}	response.response
//	@Router			/users/signup/confirm [post]
func (h *handler) confirmSignUp(c pkg.Context) {
	ctx, span := c.StartSpan()
	defer span.End()

	var verify = &confirmSignUp{}
	err := c.ParseBody(verify)
	if err != nil {
		c.Respond(response.Make(response.BadRequestCode).WithMessage("ParseBody: " + err.Error()))
		return
	}

	errs := c.ValidateStruct(verify)
	if len(errs) != 0 {
		c.Respond(response.Make(response.BadRequestCode).WithMessage("ValidateStruct: " + strings.Join(errs, "; ")))
		return
	}

	err = h.service.confirmSignUp(ctx, verify)
	if err != nil {
		span.RecordError(errors.Wrap(err, "h.service.confirmSignUp"))
		c.Respond(response.Make(response.BadGatewayCode))
		return
	}

	c.Respond(response.Make(response.OKCode))
}

// SignIn godoc
//
//	@Summary		User SignIn
//	@Tags			users-auth
//	@Description	Verifies user credentials (username/email and password) during login. If the credentials are valid, generates and returns an authentication token.
//	@Accept			json
//	@Produce		json
//	@Param			Request	body		signIn	true	"sign in info"
//	@Success		200		{object}	response.response{payload=token}
//	@Router			/users/signin [post]
func (h *handler) signIn(c pkg.Context) {
	ctx, span := c.StartSpan()
	defer span.End()

	var signIn = &signIn{}
	err := c.ParseBody(signIn)
	if err != nil {
		c.Respond(response.Make(response.BadRequestCode).WithMessage("ParseBody: " + err.Error()))
		return
	}

	errs := c.ValidateStruct(signIn)
	if len(errs) != 0 {
		c.Respond(response.Make(response.BadRequestCode).WithMessage("ValidateStruct: " + strings.Join(errs, "; ")))
		return
	}

	token, err := h.service.signIn(ctx, signIn)
	if err != nil {
		span.RecordError(errors.Wrap(err, "h.service.signIn"))
		c.Respond(response.Make(response.BadGatewayCode))
		return
	}

	c.Respond(response.Make(response.OKCode).WithPayload(token))
}

// RefreshToken godoc
//
//	@Summary		Refresh User Tokens
//	@Tags			users-auth
//	@Description	refresh user tokens
//	@Accept			json
//	@Produce		json
//	@Param			userid	path		string			true	"user id"
//	@Param			Request	body		refreshToken	true	"refresh token info"
//	@Success		200		{object}	response.response{payload=token}
//	@Router			/users/token/refresh/{userid} [post]
func (h *handler) refreshToken(c pkg.Context) {
	ctx, span := c.StartSpan()
	defer span.End()

	userID := c.GetParam("userid")
	errMessage := c.ValidateVar(userID, "required")
	if errMessage != "" {
		c.Respond(response.Make(response.BadRequestCode).WithMessage("ValidateVar: " + errMessage))
		return
	}

	refresh := &refreshToken{}
	err := c.ParseBody(refresh)
	if err != nil {
		c.Respond(response.Make(response.BadRequestCode).WithMessage("ParseBody: " + err.Error()))
		return
	}

	errs := c.ValidateStruct(refresh)
	if len(errs) != 0 {
		c.Respond(response.Make(response.BadRequestCode).WithMessage("ValidateStruct: " + strings.Join(errs, "; ")))
		return
	}

	token, err := h.service.refreshToken(ctx, userID, refresh)
	if err != nil {
		span.RecordError(errors.Wrap(err, "h.service.refreshToken"))
		c.Respond(response.Make(response.BadGatewayCode))
		return
	}

	c.Respond(response.Make(response.OKCode).WithPayload(token))
}

// GetUser godoc
//
//	@Summary		Get User Profile
//	@Tags			users-auth
//	@Description	get user profile
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	response.response{payload=user}
//	@Router			/users/profile [get]
//	@Security		Authorization Token
func (h *handler) getUser(c pkg.Context) {
	_, span := c.StartSpan()
	defer span.End()

	user, found := c.GetValue(_USER_ID_KEY).(*user)
	if !found || user == nil {
		span.RecordError(errors.New("c.GetValue"))
		c.Respond(response.Make(response.UnauthorizedCode))
		return
	}
	c.Respond(response.Make(response.OKCode).WithPayload(user))
}

// UpdateUser godoc
//
//	@Summary		Update User Profile
//	@Tags			users-auth
//	@Description	Allows users to update their profile information, such as name, contact information, or profile picture. Validates input data and updates the corresponding user record in the database.
//	@Accept			json
//	@Produce		json
//	@Param			Request	body		updateUser	true	"update user info"
//	@Success		200		{object}	response.response
//	@Router			/users/update [put]
//	@Security		Authorization Token
func (h *handler) updateUser(c pkg.Context) {
	ctx, span := c.StartSpan()
	defer span.End()

	user, found := c.GetValue(_USER_ID_KEY).(*user)
	if !found || user == nil {
		span.RecordError(errors.New("c.GetValue"))
		c.Respond(response.Make(response.UnauthorizedCode))
		return
	}

	var updateUser = &updateUser{}
	err := c.ParseBody(updateUser)
	if err != nil {
		c.Respond(response.Make(response.BadRequestCode).WithMessage("ParseBody: " + err.Error()))
		return
	}

	errs := c.ValidateStruct(updateUser)
	if len(errs) != 0 {
		c.Respond(response.Make(response.BadRequestCode).WithMessage("ValidateStruct: " + strings.Join(errs, "; ")))
		return
	}

	err = h.service.updateUser(ctx, user, updateUser)
	if err != nil {
		span.RecordError(errors.Wrap(err, "h.service.updateUser"))
		c.Respond(response.Make(response.BadGatewayCode))
		return
	}

	c.Respond(response.Make(response.OKCode))
}

// ConfirmChangeEmail godoc
//
//	@Summary		Confirm Change Email
//	@Tags			users-auth
//	@Description	Verifies a user's email address using a verification code. Activates the user's account after successful verification.
//	@Accept			json
//	@Produce		json
//	@Param			Request	body		confirmChangeEmail	true	"confirm change email info"
//	@Success		200		{object}	response.response
//	@Router			/users/update/confirm [put]
func (h *handler) confirmChangeEmail(c pkg.Context) {
	ctx, span := c.StartSpan()
	defer span.End()

	var verify = &confirmChangeEmail{}
	err := c.ParseBody(verify)
	if err != nil {
		c.Respond(response.Make(response.BadRequestCode).WithMessage("ParseBody: " + err.Error()))
		return
	}

	errs := c.ValidateStruct(verify)
	if len(errs) != 0 {
		c.Respond(response.Make(response.BadRequestCode).WithMessage("ValidateStruct: " + strings.Join(errs, "; ")))
		return
	}

	err = h.service.confirmChangeEmail(ctx, verify)
	if err != nil {
		span.RecordError(errors.Wrap(err, "h.service.confirmChangeEmail"))
		c.Respond(response.Make(response.BadGatewayCode))
		return
	}

	c.Respond(response.Make(response.OKCode))
}

// UpdateAvatar godoc
//
//	@Summary		Update User Profile Picture
//	@Tags			users-auth
//	@Description	updates user profile picture
//	@Accept			mpfd
//	@Produce		json
//	@Param			avatar	formData	file	true	"avatar file"
//	@Success		200		{object}	response.response
//	@Router			/users/update/avatar [put]
//	@Security		Authorization Token
func (h *handler) updateAvatar(c pkg.Context) {
	ctx, span := c.StartSpan()
	defer span.End()

	user, found := c.GetValue(_USER_ID_KEY).(*user)
	if !found || user == nil {
		span.RecordError(errors.New("c.GetValue"))
		c.Respond(response.Make(response.UnauthorizedCode))
		return
	}

	file, fileInfo, err := c.OpenFormFile("avatar")
	if err != nil {
		span.RecordError(errors.Wrap(err, "c.OpenFormFile"))
		c.Respond(response.Make(response.NotFoundCode).WithMessage("OpenFormFile: " + err.Error()))
		return
	}
	defer file.Close()

	err = h.service.updateAvatar(ctx, user, file, fileInfo)
	if err != nil {
		span.RecordError(errors.Wrap(err, "h.service.updateAvatar"))
		c.Respond(response.Make(response.InternalServerErrorCode))
		return
	}

	c.Respond(response.Make(response.OKCode))
}

// ForgotPassword godoc
//
//	@Summary		Forgot User Password
//	@Tags			users-auth
//	@Description	Initiates the process for resetting a forgotten password. Sends a password reset link or temporary password to the user's email, allowing them to set a new password.
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	response.response
//	@Router			/users/password/forgot [get]
//	@Security		Authorization Token
func (h *handler) forgotPassword(c pkg.Context) {
	ctx, span := c.StartSpan()
	defer span.End()

	user, found := c.GetValue(_USER_ID_KEY).(*user)
	if !found || user == nil {
		span.RecordError(errors.New("c.GetValue"))
		c.Respond(response.Make(response.UnauthorizedCode))
		return
	}

	err := h.service.forgotPassword(ctx, user.ID)
	if err != nil {
		span.RecordError(errors.Wrap(err, "h.service.forgotPassword"))
		c.Respond(response.Make(response.BadGatewayCode))
		return
	}

	c.Respond(response.Make(response.OKCode))
}

// ResetPassword godoc
//
//	@Summary		Reset User Password
//	@Tags			users-auth
//	@Description	resets user's password with confirmation code
//	@Accept			json
//	@Produce		json
//	@Param			Request	body		resetPassword	true	"reset password info"
//	@Success		200		{object}	response.response
//	@Router			/users/password/reset [put]
//	@Security		Authorization Token
func (h *handler) resetPassword(c pkg.Context) {
	ctx, span := c.StartSpan()
	defer span.End()

	var reset = &resetPassword{}
	err := c.ParseBody(reset)
	if err != nil {
		c.Respond(response.Make(response.BadRequestCode).WithMessage("ParseBody: " + err.Error()))
		return
	}

	errs := c.ValidateStruct(reset)
	if len(errs) != 0 {
		c.Respond(response.Make(response.BadRequestCode).WithMessage("ValidateStruct: " + strings.Join(errs, "; ")))
		return
	}

	err = h.service.resetPassword(ctx, reset)
	if err != nil {
		span.RecordError(errors.Wrap(err, "h.service.resetPassword"))
		c.Respond(response.Make(response.BadGatewayCode))
		return
	}

	c.Respond(response.Make(response.OKCode))
}

// ChangePassword godoc
//
//	@Summary		Change User Password
//	@Tags			users-auth
//	@Description	Enables users to change their passwords. Validates the old password, ensures password strength for the new password, and updates the password in the database.
//	@Accept			json
//	@Produce		json
//	@Param			Request	body		changePassword	true	"change password info"
//	@Success		200		{object}	response.response
//	@Router			/users/password/change [put]
//	@Security		Authorization Token
func (h *handler) changePassword(c pkg.Context) {
	ctx, span := c.StartSpan()
	defer span.End()

	user, found := c.GetValue(_USER_ID_KEY).(*user)
	if !found || user == nil {
		span.RecordError(errors.New("c.GetValue"))
		c.Respond(response.Make(response.UnauthorizedCode))
		return
	}

	var change = &changePassword{}
	err := c.ParseBody(change)
	if err != nil {
		c.Respond(response.Make(response.BadRequestCode).WithMessage("ParseBody: " + err.Error()))
		return
	}

	errs := c.ValidateStruct(change)
	if len(errs) != 0 {
		c.Respond(response.Make(response.BadRequestCode).WithMessage("ValidateStruct: " + strings.Join(errs, "; ")))
		return
	}

	err = h.service.changePassword(ctx, user, change)
	if err != nil {
		span.RecordError(errors.Wrap(err, "h.service.changePassword"))
		c.Respond(response.Make(response.BadGatewayCode))
		return
	}

	c.Respond(response.Make(response.OKCode))
}

// SignOut godoc
//
//	@Summary		User SignOut
//	@Tags			users-auth
//	@Description	log out
//	@Accept			json
//	@Produce		json
//	@Param			Request	body		refreshToken	true	"sign out info"
//	@Success		200		{object}	response.response
//	@Router			/users/signout [post]
//	@Security		Authorization Token
func (h *handler) signOut(c pkg.Context) {
	ctx, span := c.StartSpan()
	defer span.End()

	user, found := c.GetValue(_USER_ID_KEY).(*user)
	if !found || user == nil {
		span.RecordError(errors.New("c.GetValue"))
		c.Respond(response.Make(response.UnauthorizedCode))
		return
	}

	refresh := &refreshToken{}
	err := c.ParseBody(refresh)
	if err != nil {
		c.Respond(response.Make(response.BadRequestCode).WithMessage("ParseBody: " + err.Error()))
		return
	}

	errs := c.ValidateStruct(refresh)
	if len(errs) != 0 {
		c.Respond(response.Make(response.BadRequestCode).WithMessage("ValidateStruct: " + strings.Join(errs, "; ")))
		return
	}

	err = h.service.signOut(ctx, user.ID, refresh)
	if err != nil {
		span.RecordError(errors.Wrap(err, "h.service.signOut"))
		c.Respond(response.Make(response.BadGatewayCode))
		return
	}

	c.Respond(response.Make(response.OKCode))
}

// DeleteUser godoc
//
//	@Summary		Delete User Account
//	@Tags			users-auth
//	@Description	Allows users to delete their accounts. Performs necessary clean-up actions, such as revoking authentication tokens and deactivating users (not removing user data from the system).
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	response.response
//	@Router			/users/delete [delete]
//	@Security		Authorization Token
func (h *handler) deleteUser(c pkg.Context) {
	ctx, span := c.StartSpan()
	defer span.End()

	user, found := c.GetValue(_USER_ID_KEY).(*user)
	if !found || user == nil {
		span.RecordError(errors.New("c.GetValue"))
		c.Respond(response.Make(response.UnauthorizedCode))
		return
	}

	err := h.service.deleteUser(ctx, user.ID)
	if err != nil {
		span.RecordError(errors.Wrap(err, "h.service.deleteUser"))
		c.Respond(response.Make(response.BadGatewayCode))
		return
	}

	c.Respond(response.Make(response.OKCode))
}

// GetPartnerProducts godoc
//
//	@Summary		Get Partner Products
//	@Tags			partners
//	@Description	Returns available products of every partner
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	response.response{payload=products.GetAllResponse}
//	@Router			/partners/products [get]
//	@Security		Authorization Token
func (h *handler) getPartnerProducts(c pkg.Context) {
	ctx, span := c.StartSpan()
	defer span.End()

	products, err := h.service.getPartnerProducts(ctx)
	if err != nil {
		span.RecordError(errors.Wrap(err, "h.service.getPartnerProducts"))
		c.Respond(response.Make(response.InternalServerErrorCode))
		return
	}

	c.Respond(response.Make(response.OKCode).WithPayload(products))
}

// CheckOrder godoc
//
//	@Summary		Check Order
//	@Tags			orders
//	@Description	user checks an order
//	@Accept			json
//	@Produce		json
//	@Param			Request	body		checkRequest	true	"check order info"
//	@Success		200		{object}	response.response
//	@Router			/orders/check [post]
//	@Security		Authorization Token
func (h *handler) checkOrder(c pkg.Context) {
	ctx, span := c.StartSpan()
	defer span.End()

	user, found := c.GetValue(_USER_ID_KEY).(*user)
	if !found || user == nil {
		span.RecordError(errors.New("c.GetValue"))
		c.Respond(response.Make(response.UnauthorizedCode))
		return
	}

	var req = &checkRequest{}
	err := c.ParseBody(req)
	if err != nil {
		c.Respond(response.Make(response.BadRequestCode).WithMessage("ParseBody: " + err.Error()))
		return
	}

	errs := c.ValidateStruct(req)
	if len(errs) != 0 {
		c.Respond(response.Make(response.BadRequestCode).WithMessage("ValidateStruct: " + strings.Join(errs, "; ")))
		return
	}

	err = h.service.checkOrder(ctx, user, req)
	if err != nil {
		span.RecordError(errors.Wrap(err, "h.service.checkOrder"))
		c.Respond(response.Make(response.InternalServerErrorCode))
		return
	}

	c.Respond(response.Make(response.OKCode))
}

// ConfirmOrder godoc
//
//	@Summary		Confirm an Order
//	@Tags			orders
//	@Description	user confirms the order
//	@Accept			json
//	@Produce		json
//	@Param			Request	body		confirmRequest	true	"confirm order info"
//	@Success		200		{object}	response.response{payload=confirmResponse}
//	@Router			/orders/confirm [post]
//	@Security		Authorization Token
func (h *handler) confirmOrder(c pkg.Context) {
	ctx, span := c.StartSpan()
	defer span.End()

	user, found := c.GetValue(_USER_ID_KEY).(*user)
	if !found || user == nil {
		span.RecordError(errors.New("c.GetValue"))
		c.Respond(response.Make(response.UnauthorizedCode))
		return
	}

	var req = &confirmRequest{}
	err := c.ParseBody(req)
	if err != nil {
		c.Respond(response.Make(response.BadRequestCode).WithMessage("ParseBody: " + err.Error()))
		return
	}

	errs := c.ValidateStruct(req)
	if len(errs) != 0 {
		c.Respond(response.Make(response.BadRequestCode).WithMessage("ValidateStruct: " + strings.Join(errs, "; ")))
		return
	}

	resp, err := h.service.confirmOrder(ctx, user, req)
	if err != nil {
		span.RecordError(errors.Wrap(err, "h.service.confirmOrder"))
		c.Respond(response.Make(response.InternalServerErrorCode))
		return
	}

	c.Respond(response.Make(response.PendingCode).WithPayload(resp))
}

// PayOrder godoc
//
//	@Summary		Pay an Order Callback
//	@Tags			orders
//	@Description	bank sends order payment callback
//	@Accept			json
//	@Produce		json
//	@Param			Request	body		payRequest	true	"pay order info"
//	@Success		200		{object}	response.response{payload=payResponse}
//	@Router			/orders/pay [post]
func (h *handler) payOrder(c pkg.Context) {
	ctx, span := c.StartSpan()
	defer span.End()

	var req = &payRequest{}
	err := c.ParseBody(req)
	if err != nil {
		c.Respond(response.Make(response.BadRequestCode).WithMessage("ParseBody: " + err.Error()))
		return
	}

	errs := c.ValidateStruct(req)
	if len(errs) != 0 {
		c.Respond(response.Make(response.BadRequestCode).WithMessage("ValidateStruct: " + strings.Join(errs, "; ")))
		return
	}

	resp, err := h.service.payOrder(ctx, req)
	if err != nil {
		span.RecordError(errors.Wrap(err, "h.service.payOrder"))
		c.Respond(response.Make(response.InternalServerErrorCode))
		return
	}

	c.Respond(response.Make(response.OKCode).WithPayload(resp))
}

// PickUpOrder godoc
//
//	@Summary		Pick Up the Order
//	@Tags			orders
//	@Description	partner sends a callback that an order is ready
//	@Accept			json
//	@Produce		json
//	@Param			Request	body		pickupRequest	true	"pick up order info"
//	@Success		200		{object}	response.response
//	@Router			/orders/pickup [post]
func (h *handler) pickUpOrder(c pkg.Context) {
	ctx, span := c.StartSpan()
	defer span.End()

	var req = &pickupRequest{}
	err := c.ParseBody(req)
	if err != nil {
		c.Respond(response.Make(response.BadRequestCode).WithMessage("ParseBody: " + err.Error()))
		return
	}

	errs := c.ValidateStruct(req)
	if len(errs) != 0 {
		c.Respond(response.Make(response.BadRequestCode).WithMessage("ValidateStruct: " + strings.Join(errs, "; ")))
		return
	}

	err = h.service.pickUpOrder(ctx, req)
	if err != nil {
		span.RecordError(errors.Wrap(err, "h.service.pickUpOrder"))
		c.Respond(response.Make(response.InternalServerErrorCode))
		return
	}

	c.Respond(response.Make(response.OKCode))
}

// AssignOrder godoc
//
//	@Summary		Assign the Order
//	@Tags			orders
//	@Description	deliver chooses an orders and assigns it to themself
//	@Accept			json
//	@Produce		json
//	@Param			Request	body		assignRequest	true	"assign order info"
//	@Success		200		{object}	response.response
//	@Router			/orders/assign [post]
//	@Security		Authorization Token
func (h *handler) assignOrder(c pkg.Context) {
	ctx, span := c.StartSpan()
	defer span.End()

	user, found := c.GetValue(_USER_ID_KEY).(*user)
	if !found || user == nil {
		span.RecordError(errors.New("c.GetValue"))
		c.Respond(response.Make(response.UnauthorizedCode))
		return
	}

	var assignRequest = &assignRequest{}
	err := c.ParseBody(assignRequest)
	if err != nil {
		c.Respond(response.Make(response.BadRequestCode).WithMessage("ParseBody: " + err.Error()))
		return
	}

	errs := c.ValidateStruct(assignRequest)
	if len(errs) != 0 {
		c.Respond(response.Make(response.BadRequestCode).WithMessage("ValidateStruct: " + strings.Join(errs, "; ")))
		return
	}

	err = h.service.assignOrder(ctx, user, assignRequest)
	if err != nil {
		span.RecordError(errors.Wrap(err, "h.service.assignOrder"))
		c.Respond(response.Make(response.InternalServerErrorCode))
		return
	}

	c.Respond(response.Make(response.OKCode))
}

func (h *handler) authorize(actionName string, next pkg.Handler) pkg.Handler {
	return pkg.Handler(func(c pkg.Context) {
		header := c.GetHeader("Authorization") // check user role in Token and give access => authorize
		if header == "" {
			message := "empty auth header"
			statusCode := response.UnauthorizedCode
			c.Respond(response.Make(statusCode).WithMessage(message))
			return
		}

		authPrefix := "Bearer " // with a trailing space
		if !strings.Contains(header, authPrefix) {
			message := "invalid auth header"
			statusCode := response.UnauthorizedCode
			c.Respond(response.Make(statusCode).WithMessage(message))
			return
		}

		accessToken := strings.TrimPrefix(header, authPrefix)
		if len(accessToken) == 0 {
			message := "token is empty"
			statusCode := response.UnauthorizedCode
			c.Respond(response.Make(statusCode).WithMessage(message))
			return
		}

		usr, err := h.service.checkToken(c.GetRequestContext(), accessToken)
		if err != nil {
			c.Respond(response.Make(response.BadGatewayCode))
			return
		}

		// TODO: check access with actionName and usr.Roles

		user := &user{ // it should be pointer because of being comparable with nil
			ID:          usr.ID,
			FirstName:   usr.FirstName,
			LastName:    usr.LastName,
			Email:       usr.Email,
			Phone:       usr.Phone,
			PhotoURL:    usr.PhotoURL,
			BirthDate:   usr.BirthDate,
			Roles:       usr.Roles,
			NotifToken:  usr.NotifToken,
			AccessToken: accessToken,
		}

		// save user in context for further using
		c.SaveValue(_USER_ID_KEY, user)

		next(c)
	})
}

func generateHashByKey(key, value string) (string, error) {
	var h = hmac.New(sha256.New, []byte(key))
	_, err := h.Write([]byte(value))
	if err != nil {
		return "", errors.Wrap(err, "Write")
	}
	return hex.EncodeToString(h.Sum(nil)), nil
}
