package booking

import (
	"github.com/baimhons/stadiumhub/internal/match"
	"github.com/baimhons/stadiumhub/internal/seat"
	"github.com/baimhons/stadiumhub/internal/user"
	"github.com/baimhons/stadiumhub/internal/utils"
	"github.com/google/uuid"
)

// type BookingStatus string

// const (
// 	BookingPending  BookingStatus = "PENDING"
// 	BookingPaid     BookingStatus = "PAID"
// 	BookingCanceled BookingStatus = "CANCELED"
// )

type Booking struct {
	utils.BaseEntity
	UserID     uuid.UUID   `gorm:"not null"`
	User       user.User   `gorm:"foreignKey:UserID"`
	MatchID    uuid.UUID   `gorm:"not null"`
	Match      match.Match `gorm:"foreignKey:MatchID"`
	SeatID     uuid.UUID   `gorm:"not null"`
	Seat       seat.Seat   `gorm:"foreignKey:SeatID"`
	Quantity   int         `gorm:"not null"`
	TotalPrice int         `gorm:"not null"`
	Status     string      `gorm:"type:varchar(20);not null"`
}
