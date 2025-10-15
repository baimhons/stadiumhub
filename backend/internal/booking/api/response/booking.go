package response

import (
	"time"

	"github.com/google/uuid"
)

type BookingResponse struct {
	ID         uuid.UUID         `json:"id"`
	UserID     uuid.UUID         `json:"user_id"`
	MatchID    int               `json:"match_id"`
	TotalPrice int               `json:"total_price"`
	Status     string            `json:"status"`
	Seats      []BookingSeatResp `json:"seats"`
	CreatedAt  time.Time         `json:"created_at"`
	UpdatedAt  time.Time         `json:"updated_at"`
}

type BookingSeatResp struct {
	SeatID uuid.UUID `json:"seat_id"`
	Price  int       `json:"price"`
}
