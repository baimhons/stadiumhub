package request

type PaymentIntentRequest struct {
	Amount        int64  `json:"amount" validate:"required"`
	Currency      string `json:"currency" validate:"required"`
	PaymentMethod string `json:"payment_method" validate:"required"`
}
