package keycloak

import (
	"context"
	"os"
	"strings"

	"github.com/Nerzal/gocloak/v13"
	"github.com/pkg/errors"
	"github.com/shahzodshafizod/gocloud/pkg"
)

type auth struct {
	realm        string
	clientID     string
	clientSecret string
	client       *gocloak.GoCloak
}

func NewAuth() pkg.Auth {
	return &auth{
		realm:        os.Getenv("KEYCLOAK_REALM"),
		clientID:     os.Getenv("KEYCLOAK_CLIENT_ID"),
		clientSecret: os.Getenv("KEYCLOAK_CLIENT_SECRET"),
		client:       gocloak.NewClient(os.Getenv("KEYCLOAK_BASE_URL")),
	}
}

func (a *auth) SignUp(ctx context.Context, sr *pkg.SignUp) (string, error) {
	// 1. get admin token to create the user
	adminToken, err := a.client.LoginClient(ctx, a.clientID, a.clientSecret, a.realm)
	if err != nil {
		return "", errors.Wrap(err, "a.client.LoginClient")
	}

	// 2. create user
	user := gocloak.User{
		FirstName:     gocloak.StringP(sr.FirstName),
		LastName:      gocloak.StringP(sr.LastName),
		Email:         gocloak.StringP(sr.Email),
		EmailVerified: gocloak.BoolP(false),
		Enabled:       gocloak.BoolP(true),
	}

	attributes := make(map[string][]string)
	if sr.Phone != "" {
		attributes["phoneNumber"] = []string{sr.Phone}
	}
	if sr.BirthDate != "" {
		attributes["birthdate"] = []string{sr.BirthDate}
	}
	if sr.NotifToken != "" {
		attributes["notif_token"] = []string{sr.NotifToken}
	}

	user.Attributes = &attributes

	userID, err := a.client.CreateUser(ctx, adminToken.AccessToken, a.realm, user)
	if err != nil {
		return "", errors.Wrap(err, "a.client.CreateUser")
	}

	// 3. set password
	err = a.client.SetPassword(ctx, adminToken.AccessToken, userID, a.realm, sr.Password, false)
	if err != nil {
		return "", errors.Wrap(err, "a.client.SetPassword")
	}

	// 4. add role to user
	sr.Role = strings.ToLower(sr.Role)
	role, err := a.client.GetRealmRole(ctx, adminToken.AccessToken, a.realm, sr.Role)
	if err != nil {
		return "", errors.Wrap(err, "a.client.GetRealmRole")
	}
	err = a.client.AddRealmRoleToUser(ctx, adminToken.AccessToken, a.realm, userID, []gocloak.Role{*role})
	if err != nil {
		return "", errors.Wrap(err, "a.client.AddRealmRoleToUser")
	}

	// 5. send email verification link
	err = a.client.SendVerifyEmail(ctx, adminToken.AccessToken, userID, a.realm)
	if err != nil {
		return "", errors.Wrap(err, "a.client.SendVerifyEmail")
	}

	return userID, nil
}

func (a *auth) ConfirmSignUp(ctx context.Context, ve *pkg.VerifyEmail) error {
	return nil
}

func (a *auth) SignIn(ctx context.Context, signIn *pkg.SignIn) (*pkg.User, *pkg.Token, error) {
	userToken, err := a.client.Login(ctx, a.clientID, a.clientSecret, a.realm, signIn.Email, signIn.Password)
	if err != nil {
		return nil, nil, errors.Wrap(err, "a.client.Login")
	}

	user, err := a.getUser(ctx, userToken.AccessToken)
	if err != nil {
		return nil, nil, errors.Wrap(err, "a.getUser")
	}

	token := &pkg.Token{
		AccessToken:  userToken.AccessToken,
		ExpiresIn:    userToken.ExpiresIn,
		RefreshToken: userToken.RefreshToken,
	}

	return user, token, nil
}

func (a *auth) getUser(ctx context.Context, accessToken string) (*pkg.User, error) {
	info, err := a.client.GetRawUserInfo(ctx, accessToken, a.realm)
	if err != nil {
		return nil, errors.Wrap(err, "a.client.GetRawUserInfo")
	}

	// it's more clearly to use just GetUserInfo method to get user, but that doesn't have
	// some attributes, e.g. user roles. Decoding access token is for missing fields
	_, claims, err := a.client.DecodeAccessToken(ctx, accessToken, a.realm)
	if err != nil {
		return nil, errors.Wrap(err, "a.client.DecodeAccessToken")
	}

	roles := make([]string, 0)
	for _, value := range parseAny[[]any]((parseAny[map[string]any]((*claims)["realm_access"]))["roles"]) {
		roles = append(roles, parseAny[string](value))
	}

	return &pkg.User{
		ID:         parseAny[string](info["sub"]),
		FirstName:  parseAny[string](info["given_name"]),
		LastName:   parseAny[string](info["family_name"]),
		Email:      parseAny[string](info["email"]),
		Phone:      parseAny[string](info["phone_number"]),
		PhotoURL:   parseAny[string](info["picture"]),
		BirthDate:  parseAny[string](info["birthdate"]),
		Roles:      roles,
		NotifToken: parseAny[string](info["notif_token"]),
	}, nil
}

func (a *auth) CheckToken(ctx context.Context, accessToken string) (*pkg.User, error) {
	result, err := a.client.RetrospectToken(ctx, accessToken, a.clientID, a.clientSecret, a.realm)
	if err != nil {
		return nil, errors.Wrap(err, "a.client.RetrospectToken")
	}
	if !gocloak.PBool(result.Active) {
		return nil, errors.New("token is invalid")
	}

	return a.getUser(ctx, accessToken)
}

func (a *auth) RefreshToken(ctx context.Context, _ string, refreshToken string) (*pkg.Token, error) {
	token, err := a.client.RefreshToken(ctx, refreshToken, a.clientID, a.clientSecret, a.realm)
	if err != nil {
		return nil, errors.Wrap(err, "a.client.RefreshToken")
	}
	return &pkg.Token{
		AccessToken:  token.AccessToken,
		RefreshToken: refreshToken,
		ExpiresIn:    token.ExpiresIn,
	}, nil
}

func (a *auth) UpdateUser(ctx context.Context, _ string, updateUser *pkg.UpdateUser) error {
	// 1. get admin token for updating the user
	adminToken, err := a.client.LoginClient(ctx, a.clientID, a.clientSecret, a.realm)
	if err != nil {
		return errors.Wrap(err, "a.client.LoginClient")
	}

	// 2. get user to set new fields to update
	user, err := a.client.GetUserByID(ctx, adminToken.AccessToken, a.realm, updateUser.ID)
	if err != nil {
		return errors.Wrap(err, "a.client.GetUserByID")
	}

	if user.Attributes == nil {
		user.Attributes = &map[string][]string{}
	}

	// 3. set user fields to update
	if updateUser.FirstName != nil {
		user.FirstName = updateUser.FirstName
	}
	if updateUser.LastName != nil {
		user.LastName = updateUser.LastName
	}
	var emailChanged = false
	if updateUser.Email != nil {
		if gocloak.PString(user.Email) != gocloak.PString(updateUser.Email) {
			emailChanged = true
			user.EmailVerified = gocloak.BoolP(false)
			user.Email = updateUser.Email
		}
	}
	if updateUser.Phone != nil {
		(*user.Attributes)["phoneNumber"] = []string{gocloak.PString(updateUser.Phone)}
	}
	if updateUser.PhotoURL != nil {
		(*user.Attributes)["picture"] = []string{gocloak.PString(updateUser.PhotoURL)}
	}
	if updateUser.BirthDate != nil {
		(*user.Attributes)["birthdate"] = []string{gocloak.PString(updateUser.BirthDate)}
	}

	// 4. update the user
	err = a.client.UpdateUser(ctx, adminToken.AccessToken, a.realm, *user)
	if err != nil {
		return errors.Wrap(err, "a.client.UpdateUser")
	}

	// 5. send email verification
	if emailChanged {
		err = a.client.LogoutAllSessions(ctx, adminToken.AccessToken, a.realm, updateUser.ID)
		if err != nil {
			return errors.Wrap(err, "a.client.LogoutAllSessions")
		}

		err = a.client.SendVerifyEmail(ctx, adminToken.AccessToken, updateUser.ID, a.realm)
		if err != nil {
			return errors.Wrap(err, "a.client.SendVerifyEmail")
		}
	}

	return nil
}

func (a *auth) ConfirmChangeEmail(ctx context.Context, accesstoken string, ve *pkg.VerifyEmail) error {
	return nil
}

func (a *auth) ForgotPassword(ctx context.Context, userID string) error {
	adminToken, err := a.client.LoginClient(ctx, a.clientID, a.clientSecret, a.realm)
	if err != nil {
		return errors.Wrap(err, "a.client.LoginClient")
	}

	return a.client.ExecuteActionsEmail(ctx, adminToken.AccessToken, a.realm, gocloak.ExecuteActionsEmail{
		UserID:   gocloak.StringP(userID),
		ClientID: gocloak.StringP(a.clientID),
		Actions:  &([]string{"UPDATE_PASSWORD"}),
	})
}

func (a *auth) ResetPassword(ctx context.Context, reset *pkg.ResetPassword) error {
	return nil
}

func (a *auth) ChangePassword(ctx context.Context, passwords *pkg.ChangePassword) error {
	// 1. log in to check the old password
	_, err := a.client.Login(ctx, a.clientID, a.clientSecret, a.realm, passwords.Email, passwords.OldPassword)
	if err != nil {
		return errors.Wrap(err, "a.client.Login")
	}

	// 2. get admin token to check old password and set new password
	adminToken, err := a.client.LoginClient(ctx, a.clientID, a.clientSecret, a.realm)
	if err != nil {
		return errors.Wrap(err, "a.client.LoginClient")
	}

	// 3. password will be reset, so all released token should be revoked
	err = a.client.LogoutAllSessions(ctx, adminToken.AccessToken,
		a.realm, gocloak.PString(&passwords.UserID))
	if err != nil {
		return errors.Wrap(err, "a.client.LogoutAllSessions")
	}

	// 4. change password
	return a.client.SetPassword(ctx, adminToken.AccessToken,
		passwords.UserID, a.realm, passwords.NewPassword, false)
}

func (a *auth) SignOut(ctx context.Context, _ string, refreshToken string) error {
	return a.client.Logout(ctx, a.clientID, a.clientSecret, a.realm, refreshToken)
}

func (a *auth) DeleteUser(ctx context.Context, userID string) error {
	adminToken, err := a.client.LoginClient(ctx, a.clientID, a.clientSecret, a.realm)
	if err != nil {
		return errors.Wrap(err, "a.client.LoginClient")
	}

	// ! Do NOT delete users, just inactive them.
	// a.client.DeleteUser(ctx, token.AccessToken, a.realm, userID)

	return a.client.UpdateUser(ctx, adminToken.AccessToken, a.realm, gocloak.User{
		ID:      gocloak.StringP(userID),
		Enabled: gocloak.BoolP(false),
	})
}

func parseAny[T any](v any) (value T) {
	if v == nil {
		return
	}
	value, _ = v.(T)
	return
}
