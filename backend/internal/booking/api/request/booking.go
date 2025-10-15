package request

import (
	"github.com/google/uuid"
)

type CreateBookingRequest struct {
	MatchID int         `json:"match_id" binding:"required"`
	SeatIDs []uuid.UUID `json:"seat_ids" binding:"required,min=1"`
}

type CancelBookingRequest struct {
	BookingID uuid.UUID `json:"booking_id" binding:"required"`
}

type UpdateBookingStatusRequest struct {
	BookingID uuid.UUID `json:"booking_id" binding:"required"`
	Status    string    `json:"status" binding:"required,oneof=PENDING PAID CANCELED"`
}
