package request

import "github.com/google/uuid"

type PaymentRequest struct {
	BookingID     uuid.UUID `json:"booking_id" binding:"required,uuid"`
	PaymentMethod string    `json:"payment_method" binding:"required"`
	Proof         string    `json:"proof"`
}
