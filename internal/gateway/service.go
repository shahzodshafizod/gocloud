package gateway

import (
	"bytes"
	"context"
	"encoding/json"
	"image"
	"image/jpeg"
	"image/png"
	"io"
	"time"

	"github.com/disintegration/imaging"
	gonanoid "github.com/matoous/go-nanoid/v2"
	"github.com/pkg/errors"
	"github.com/shahzodshafizod/gocloud/internal/orders"
	"github.com/shahzodshafizod/gocloud/internal/partners"
	"github.com/shahzodshafizod/gocloud/internal/products"
	"github.com/shahzodshafizod/gocloud/pkg"
)

type Service interface {
	signUp(context.Context, *signUp) (*signUpResponse, error)
	confirmSignUp(context.Context, *confirmSignUp) error
	signIn(context.Context, *signIn) (*token, error)
	checkToken(context.Context, string) (*pkg.User, error)
	refreshToken(context.Context, string, *refreshToken) (*token, error)
	updateUser(context.Context, *user, *updateUser) error
	confirmChangeEmail(context.Context, *confirmChangeEmail) error
	updateAvatar(context.Context, *user, pkg.File, pkg.FileInfo) error
	forgotPassword(context.Context, string) error
	resetPassword(context.Context, *resetPassword) error
	changePassword(context.Context, *user, *changePassword) error
	signOut(context.Context, string, *refreshToken) error
	deleteUser(context.Context, string) error

	getPartnerProducts(context.Context) (*products.GetAllResponse, error)
	checkOrder(context.Context, *user, *checkRequest) error
	confirmOrder(context.Context, *user, *confirmRequest) (*confirmResponse, error)
	payOrder(context.Context, *payRequest) (*payResponse, error)
	pickUpOrder(context.Context, *pickupRequest) error
	assignOrder(context.Context, *user, *assignRequest) error
}

type service struct {
	authManager pkg.Auth
	cache       pkg.Cache
	queue       pkg.Publisher
	partners    partners.PartnersClient
	orders      orders.OrdersClient
	storage     pkg.Storage
}

func NewService(
	authManager pkg.Auth,
	partnersClient partners.PartnersClient,
	ordersClient orders.OrdersClient,
	cache pkg.Cache,
	queue pkg.Queue,
	storage pkg.Storage,
) Service {
	return &service{
		authManager: authManager,
		cache:       cache,
		queue:       queue,
		storage:     storage,
		partners:    partnersClient,
		orders:      ordersClient,
	}
}

func (s service) signUp(ctx context.Context, signUp *signUp) (*signUpResponse, error) {
	userID, err := s.authManager.SignUp(ctx, &pkg.SignUp{
		FirstName:  signUp.FirstName,
		LastName:   signUp.LastName,
		Email:      signUp.Email,
		Password:   signUp.Password,
		Phone:      signUp.Phone,
		BirthDate:  signUp.BirthDate,
		Role:       signUp.Role,
		NotifToken: signUp.NotifToken,
	})
	if err != nil {
		return nil, errors.Wrap(err, "s.authManager.SignUp")
	}

	return &signUpResponse{userID}, nil
}

func (s *service) confirmSignUp(ctx context.Context, verify *confirmSignUp) error {
	err := s.authManager.ConfirmSignUp(ctx, &pkg.VerifyEmail{
		UserID: verify.UserID,
		Code:   verify.Code,
	})
	if err != nil {
		return errors.Wrap(err, "s.authManager.ConfirmSignUp")
	}
	return nil
}

func (s *service) signIn(ctx context.Context, signIn *signIn) (*token, error) {
	_, authToken, err := s.authManager.SignIn(ctx, &pkg.SignIn{
		Email:    signIn.Email,
		Password: signIn.Password,
	})
	if err != nil {
		return nil, errors.Wrap(err, "s.authManager.SignIn")
	}
	return &token{
		AccessToken:  authToken.AccessToken,
		RefreshToken: authToken.RefreshToken,
		ExpiresIn:    authToken.ExpiresIn,
	}, nil
}

func (s *service) checkToken(ctx context.Context, accessToken string) (*pkg.User, error) {
	user, err := s.authManager.CheckToken(ctx, accessToken)
	if err != nil {
		return nil, errors.Wrap(err, "s.authManager.CheckToken")
	}
	return user, nil
}

func (s *service) refreshToken(ctx context.Context, userID string, refresh *refreshToken) (*token, error) {
	authToken, err := s.authManager.RefreshToken(ctx, userID, refresh.RefreshToken)
	if err != nil {
		return nil, errors.Wrap(err, "s.authManager.RefreshToken")
	}
	return &token{
		AccessToken:  authToken.AccessToken,
		RefreshToken: authToken.RefreshToken,
		ExpiresIn:    authToken.ExpiresIn,
	}, nil
}

func (s *service) updateUser(ctx context.Context, user *user, updateUser *updateUser) error {
	userToUpdate := &pkg.UpdateUser{ID: user.ID}

	if updateUser.FirstName != "" {
		userToUpdate.FirstName = &updateUser.FirstName
	}
	if updateUser.LastName != "" {
		userToUpdate.LastName = &updateUser.LastName
	}
	if updateUser.Email != "" {
		userToUpdate.Email = &updateUser.Email
	}
	if updateUser.Phone != "" {
		userToUpdate.Phone = &updateUser.Phone
	}
	if updateUser.BirthDate != "" {
		userToUpdate.BirthDate = &updateUser.BirthDate
	}
	if updateUser.NotifToken != "" {
		userToUpdate.NotifToken = &updateUser.NotifToken
	}

	err := s.authManager.UpdateUser(ctx, user.AccessToken, userToUpdate)
	if err != nil {
		return errors.Wrap(err, "s.authManager.UpdateUser")
	}

	return nil
}

func (s *service) confirmChangeEmail(ctx context.Context, verify *confirmChangeEmail) error {
	err := s.authManager.ConfirmChangeEmail(ctx, verify.AccessToken, &pkg.VerifyEmail{
		UserID: verify.UserID,
		Code:   verify.Code,
	})
	if err != nil {
		return errors.Wrap(err, "s.authManager.ConfirmChangeEmail")
	}
	return nil
}

func (s *service) updateAvatar(ctx context.Context, user *user, file pkg.File, fileInfo pkg.FileInfo) error {
	oldPhotoURL := user.PhotoURL

	img, format, err := image.Decode(io.LimitReader(file, _MAX_AVATAR_BYTES))
	if err != nil {
		return errors.Wrap(err, "image.Decode")
	}

	img = imaging.Fill(img, 400, 400, imaging.Center, imaging.CatmullRom)

	fileName, err := gonanoid.New(20)
	if err != nil {
		return errors.Wrap(err, "gonanoid.New")
	}

	buffer := new(bytes.Buffer)
	switch format {
	case "png":
		fileName += ".png"
		err = png.Encode(buffer, img)
	case "jpeg":
		fileName += ".jpg"
		err = jpeg.Encode(buffer, img, nil)
	default:
		return errors.New("UnsupportedAvatarFormat: " + format)
	}

	if err != nil {
		return errors.Wrap(err, "format.Encode")
	}

	data := buffer.Bytes()

	photoURL, err := s.storage.Upload(ctx, pkg.UploadInput{
		File:        bytes.NewReader(data),
		Name:        fileName,
		Size:        int64(len(data)),
		ContentType: fileInfo.ContentType(),
	})
	if err != nil {
		return errors.Wrap(err, "s.storage.Upload")
	}

	err = s.authManager.UpdateUser(ctx, user.AccessToken, &pkg.UpdateUser{
		ID:       user.ID,
		PhotoURL: &photoURL,
	})
	if err != nil {
		err = errors.Wrap(err, "s.authManager.UpdateUser")
		err2 := s.storage.Delete(ctx, photoURL)
		if err2 != nil {
			err2 = errors.Wrap(err2, "s.storage.Delete")
			err = errors.Wrap(err, err2.Error())
		}
		return err
	}

	if oldPhotoURL != "" {
		err = s.storage.Delete(ctx, oldPhotoURL)
		if err != nil {
			return errors.Wrap(err, "s.storage.Delete")
		}
	}

	return nil
}

func (s *service) forgotPassword(ctx context.Context, userID string) error {
	err := s.authManager.ForgotPassword(ctx, userID)
	if err != nil {
		return errors.Wrap(err, "s.authManager.ForgotPassword")
	}
	return nil
}

func (s *service) resetPassword(ctx context.Context, reset *resetPassword) error {
	err := s.authManager.ResetPassword(ctx, &pkg.ResetPassword{
		UserID:   reset.UserID,
		Code:     reset.Code,
		Password: reset.Password,
	})
	if err != nil {
		return errors.Wrap(err, "s.authManager.ResetPassword")
	}
	return nil
}

func (s *service) changePassword(ctx context.Context, user *user, change *changePassword) error {
	err := s.authManager.ChangePassword(ctx, &pkg.ChangePassword{
		UserID:      user.ID,
		AccessToken: user.AccessToken,
		Email:       user.Email,
		OldPassword: change.OldPassword,
		NewPassword: change.NewPassword,
	})
	if err != nil {
		return errors.Wrap(err, "s.authManager.ChangePassword")
	}
	return nil
}

func (s *service) signOut(ctx context.Context, userID string, refresh *refreshToken) error {
	err := s.authManager.SignOut(ctx, userID, refresh.RefreshToken)
	if err != nil {
		return errors.Wrap(err, "s.authManager.SignOut")
	}
	return nil
}

func (s *service) deleteUser(ctx context.Context, userID string) error {
	err := s.authManager.DeleteUser(ctx, userID)
	if err != nil {
		return errors.Wrap(err, "s.authManager.DeleteUser")
	}
	return nil
}

func (s *service) getPartnerProducts(ctx context.Context) (*products.GetAllResponse, error) {
	products, err := s.partners.GetPartnerProducts(ctx, &products.GetAllRequest{})
	if err != nil {
		return nil, errors.Wrap(err, "s.partners.GetPartnerProducts")
	}
	return products, nil
}

func (s *service) checkOrder(ctx context.Context, user *user, req *checkRequest) error {
	// TODO: Ensure no duplicate OrderID exists in the Orders Service before processing the request.

	cacheKey := "ORDER::" + user.ID + req.OrderID
	var details = &orders.Order{}
	err := s.cache.GetStruct(ctx, cacheKey, details)
	if err == nil {
		return nil
	}

	// check the partner and products for availability
	checkReq := &partners.CheckRequest{
		PartnerID:   int32(req.PartnerID),
		TotalAmount: req.TotalAmount,
		Products:    make([]*orders.Product, len(req.Products)),
	}
	for idx, product := range req.Products {
		checkReq.Products[idx] = &orders.Product{
			ID:       int32(product.ID),
			Quantity: int32(product.Quantity),
		}
	}
	checkResp, err := s.partners.CheckPartnerProducts(ctx, checkReq)
	if err != nil {
		return errors.Wrap(err, "s.partners.CheckPartnerProducts")
	}

	// save order details in cache and publish in queue for orders service
	details.OrderID = req.OrderID
	details.CustomerID = user.ID
	details.CustomerName = user.FirstName + " " + user.LastName
	details.CustomerPhone = req.CustomerPhone
	details.CustomerNotifToken = user.NotifToken
	details.DeliveryAddress = req.DeliveryAddress
	details.PartnerID = req.PartnerID
	details.PartnerTitle = checkResp.PartnerTitle
	details.PartnerBrand = checkResp.PartnerBrand
	details.Products = checkResp.Products
	details.TotalAmount = req.TotalAmount
	details.Paytype = req.Paytype

	err = s.cache.SaveStruct(ctx, cacheKey, details, time.Minute*10)
	if err != nil {
		return errors.Wrap(err, "s.cache.SaveStruct")
	}

	return nil
}

func (s *service) confirmOrder(ctx context.Context, user *user, req *confirmRequest) (*confirmResponse, error) {
	var key = "ORDER::" + user.ID + req.OrderID
	var order = &orders.Order{}
	err := s.cache.GetStruct(ctx, key, order)
	if err != nil {
		return nil, errors.Wrap(err, "s.cache.GetStruct")
	}

	resp, err := s.orders.CreateOrder(ctx, order)
	if err != nil {
		return nil, errors.Wrap(err, "s.orders.CreateOrder")
	}

	// Don't worry if there is an error, bc cached data will be invalidated in 10 minutes:)
	s.cache.Del(ctx, key)

	return &confirmResponse{
		OrderID:        resp.OrderID,
		TotalAmount:    order.TotalAmount,
		PartnerTitle:   order.PartnerTitle,
		PartnerBrand:   order.PartnerBrand,
		WebcheckoutURL: resp.WebcheckoutURL,
		CallbackURL:    resp.CallbackURL,
	}, nil
}

func (s *service) payOrder(ctx context.Context, req *payRequest) (*payResponse, error) {
	resp, err := s.orders.PayOrder(ctx, &orders.PayRequest{
		OrderID:    req.OrderID,
		PaymentID:  req.PaymentID,
		PaidAmount: req.PaidAmount,
	})
	if err != nil {
		return nil, errors.Wrap(err, "s.orders.PayOrder")
	}

	return &payResponse{
		PaymentID: resp.PaymentID,
	}, nil
}

func (s *service) pickUpOrder(ctx context.Context, req *pickupRequest) error {
	msgBody, err := json.Marshal(req)
	if err != nil {
		return errors.Wrap(err, "json.Marshal")
	}
	err = s.queue.Publish(ctx, "orders.ready", msgBody)
	if err != nil {
		return errors.Wrap(err, "s.queue.Publish")
	}
	return nil
}

func (s *service) assignOrder(ctx context.Context, user *user, req *assignRequest) error {
	_, err := s.orders.AssignOrder(ctx, &orders.AssignRequest{
		OrderID:     req.OrderID,
		DelivererID: user.ID,
	})
	if err != nil {
		return errors.Wrap(err, "s.orders.AssignOrder")
	}
	return nil
}
