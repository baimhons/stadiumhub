package response

import "github.com/google/uuid"

type PaymentResponse struct {
	ID            uuid.UUID `json:"id"`
	Amount        float64   `json:"amount"`
	PaymentMethod string    `json:"payment_method"`
	Status        string    `json:"status"`
	PaymentAt     int64     `json:"payment_at"`
}
