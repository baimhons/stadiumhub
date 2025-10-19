package request

import "github.com/google/uuid"

type PaymentIntentRequest struct {
	Amount    int64     `json:"amount" validate:"required"`
	Currency  string    `json:"currency" validate:"required"`
	BookingID uuid.UUID `json:"booking_id" validate:"required"`
}
