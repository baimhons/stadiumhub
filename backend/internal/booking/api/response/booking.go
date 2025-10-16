package response

import (
	"time"

	"github.com/baimhons/stadiumhub/internal/match"
	"github.com/baimhons/stadiumhub/internal/user"
	"github.com/google/uuid"
)

type BookingResponse struct {
	ID         uuid.UUID         `json:"id"`
	UserID     uuid.UUID         `json:"user_id"`
	User       user.User         `gorm:"foreignKey:UserID"`
	MatchID    int               `json:"match_id"`
	Match      match.Match       `gorm:"foreignKey:MatchID"`
	TotalPrice int               `json:"total_price"`
	Status     string            `json:"status"`
	Seats      []BookingSeatResp `json:"seats"`
	CreatedAt  time.Time         `json:"created_at"`
	UpdatedAt  time.Time         `json:"updated_at"`
}

type BookingSeatResp struct {
	SeatID uuid.UUID `json:"seat_id"`
	SeatNo string    `json:"seat_no"`
	Price  int       `json:"price"`
}
