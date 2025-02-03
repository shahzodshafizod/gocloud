package orders

type bank struct {
	ID             string `json:"id"`
	WebcheckoutURL string `json:"webcheckout_url"`
}

type PaidOrder struct {
	OrderID     int64      `json:"order_id"`
	PartnerID   int        `json:"partner_id"`
	CallbackURL string     `json:"callback_url"`
	Products    []*Product `json:"products"`
}

type pickupRequest struct {
	OrderID       int64  `json:"order_id"`
	PickupAddress string `json:"pickup_address"` // Address from which the package will be picked up.
}

// type orderAddress struct {
// 	OrderID     int64  `json:"order_id"`
// 	Street      string `json:"street"`
// 	City        string `json:"city"`
// 	State       string `json:"state"`
// 	PostalCode  string `json:"postal_code"`
// 	AddressType string `json:"address_type"` // Billing or Shipping
// }
