package gateway

// registerUser represents the details of a user for registration.
type signUp struct {
	FirstName  string `json:"first_name" validate:"required,alpha,min=2"` // The first name of the user.
	LastName   string `json:"last_name" validate:"required,alpha,min=2"`  // The last name of the user.
	Email      string `json:"email" validate:"required,email"`            // The email address of the user.
	Password   string `json:"password" validate:"required,min=6"`         // The hashed password of the user. Note: It's important to store passwords securely by hashing them.
	Phone      string `json:"phone" validate:"omitempty,e164"`            // optional: E.164 formatted phone number: [+] [country code] [subscriber number including area code] and can have a maximum of fifteen digits.
	BirthDate  string `json:"birth_date" validate:"omitempty,dateonly"`   // optional: YYYY-MM-DD
	Role       string `json:"role" validate:"required,oneof=client deliver partner admin"`
	NotifToken string `json:"notif_token" validate:"omitempty"` // Notification Token ID
}

type user struct {
	ID          string   `json:"id"`         // A unique identifier for the user.
	FirstName   string   `json:"first_name"` // The first name of the user.
	LastName    string   `json:"last_name"`  // The last name of the user.
	Email       string   `json:"email"`      // The email address of the user.
	Phone       string   `json:"phone"`
	PhotoURL    string   `json:"photo_url"`
	BirthDate   string   `json:"birth_date"` // YYYY-MM-DD
	Roles       []string `json:"roles"`
	NotifToken  string   `json:"-"`
	AccessToken string   `json:"-"`
}

type signUpResponse struct {
	UserID string `json:"user_id"`
}

type confirmSignUp struct {
	UserID string `json:"user_id" validate:"required"`
	Code   string `json:"code" validate:"required,number"`
}

type confirmChangeEmail struct {
	UserID      string `json:"user_id" validate:"required"`
	AccessToken string `json:"access_token" validate:"required"`
	Code        string `json:"code" validate:"required,number"`
}

// Login represents the credentials used for user authentication.
type signIn struct {
	Email    string `json:"email" validate:"required,email"`    // The email address provided during login.
	Password string `json:"password" validate:"required,min=6"` // The password provided during login.
}

type token struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int    `json:"expires_in"`
}

type refreshToken struct {
	RefreshToken string `json:"refresh_token" validate:"required"`
}

// UpdateUser represents the details used for updating a user's profile.
type updateUser struct {
	FirstName  string `json:"first_name" validate:"omitempty,alpha,min=2"` // The updated first name of the user.
	LastName   string `json:"last_name" validate:"omitempty,alpha,min=2"`  // The updated last name of the user.
	Email      string `json:"email" validate:"omitempty,email"`            // The updated email address of the user.
	Phone      string `json:"phone" validate:"omitempty,e164"`             // E.164 formatted phone number: [+] [country code] [subscriber number including area code] and can have a maximum of fifteen digits.
	BirthDate  string `json:"birth_date" validate:"omitempty,dateonly"`    // YYYY-MM-DD
	NotifToken string `json:"notif_token" validate:"omitempty"`            // Notification Token ID
}

type resetPassword struct {
	UserID   string `json:"user_id" validate:"required"`
	Code     string `json:"code" validate:"required,number"`
	Password string `json:"password" validate:"required,min=6"`
}

// ChangePassword represents the request structure for changing a user's password.
type changePassword struct {
	OldPassword string `json:"old_password" validate:"required,min=6"` // The user's current password.
	NewPassword string `json:"new_password" validate:"required,min=6"` // The new password the user wants to set.
}

type checkRequest struct {
	OrderID         string     `json:"order_id" validate:"required"`
	PartnerID       int32      `json:"partner_id" validate:"required,gt=0"`
	CustomerPhone   string     `json:"customer_phone" validate:"omitempty,e164"`
	DeliveryAddress string     `json:"delivery_address" validate:"required,min=6"`
	Products        []*product `json:"products" validate:"required"`
	TotalAmount     int64      `json:"total_amount" validate:"required,gt=0"`
	Paytype         string     `json:"paytype" validate:"required"`
}

type confirmRequest struct {
	OrderID string `json:"order_id" validate:"required"`
}

type confirmResponse struct {
	OrderID        int64  `json:"order_id"`
	TotalAmount    int64  `json:"total_amount"`
	PartnerTitle   string `json:"partner_title"`
	PartnerBrand   string `json:"partner_brand"`
	WebcheckoutURL string `json:"webcheckout_url"`
	CallbackURL    string `json:"callback_url"`
}

type payRequest struct {
	OrderID    int64  `json:"order_id" validate:"required"`
	PaymentID  string `json:"payment_id" validate:"required"`
	PaidAmount int64  `json:"paid_amount" validate:"required"`
}

type payResponse struct {
	PaymentID string `json:"payment_id"`
}

type product struct {
	ID       int `json:"id" validate:"required"`
	Quantity int `json:"quantity" validate:"required"`
}

type pickupRequest struct {
	OrderID       int64  `json:"order_id"`
	PickupAddress string `json:"pickup_address" validate:"required"`
}

type assignRequest struct {
	OrderID int64 `json:"order_id"`
}
