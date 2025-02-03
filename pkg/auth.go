//go:generate mockgen -source=auth.go -package=mocks -destination=mocks/auth.go
package pkg

import (
	"context"
)

type Auth interface {
	// Registers a new user and returns a user ID or error if the process fails.
	SignUp(ctx context.Context, signUp *SignUp) (string, error)
	// Verifies a user's email address using a confirmation code.
	ConfirmSignUp(ctx context.Context, verify *VerifyEmail) error
	// Authenticates a user and returns user details, tokens, or an error.
	SignIn(ctx context.Context, signIn *SignIn) (*User, *Token, error)
	// Validates an access token and returns the associated user details.
	CheckToken(ctx context.Context, accessToken string) (*User, error)
	// Generates a new access token using a refresh token for an authenticated user.
	RefreshToken(ctx context.Context, userID string, refreshToken string) (*Token, error)
	// Updates user information (e.g., profile details) using an access token.
	UpdateUser(ctx context.Context, accessToken string, user *UpdateUser) error // accessToken is just for cognito
	// Confirms a user's email change request with a verification code.
	ConfirmChangeEmail(ctx context.Context, accessToken string, verifyEmail *VerifyEmail) error // accessToken is just for cognito
	// Initiates a password reset process for a user by their ID.
	ForgotPassword(ctx context.Context, userID string) error
	// Resets a user's password using a reset token or code.
	ResetPassword(ctx context.Context, reset *ResetPassword) error
	// Allows a user to change their password after providing the current one.
	ChangePassword(ctx context.Context, change *ChangePassword) error
	// Logs out a user by invalidating their refresh token.
	SignOut(ctx context.Context, userID string, refreshToken string) error
	// Deletes a user account by their ID.
	DeleteUser(ctx context.Context, userID string) error
}

type SignUp struct {
	FirstName  string `bson:"first_name"`
	LastName   string `bson:"last_name"`
	Email      string `bson:"email"`
	Password   string `bson:"password"`
	Phone      string `bson:"phone"`
	BirthDate  string `bson:"birth_date"`
	Role       string `bson:"role"`
	NotifToken string `bson:"notif_token"`
}

type VerifyEmail struct {
	UserID string
	Code   string
}

type SignIn struct {
	Email    string
	Password string
}

type User struct {
	ID         string
	FirstName  string
	LastName   string
	Email      string
	Phone      string
	PhotoURL   string
	BirthDate  string
	Roles      []string
	NotifToken string
}

type Token struct {
	AccessToken  string
	RefreshToken string
	ExpiresIn    int
}

type UpdateUser struct {
	ID         string
	FirstName  *string
	LastName   *string
	Email      *string
	Phone      *string
	PhotoURL   *string
	BirthDate  *string
	NotifToken *string
}

type ResetPassword struct {
	UserID   string
	Code     string
	Password string
}

type ChangePassword struct {
	UserID      string
	AccessToken string // for cognito
	Email       string
	OldPassword string
	NewPassword string
}
