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
	UserID     uuid.UUID     `gorm:"type:char(36);not null"`
	User       user.User     `gorm:"foreignKey:UserID"`
	MatchID    int           `gorm:"not null"`
	Match      match.Match   `gorm:"foreignKey:MatchID"`
	TotalPrice int           `gorm:"not null"`
	Status     string        `gorm:"type:varchar(20);not null"`
	Seats      []BookingSeat `gorm:"foreignKey:BookingID"`
}

type BookingSeat struct {
	utils.BaseEntity
	BookingID uuid.UUID `gorm:"type:char(36);not null"`
	Booking   Booking   `gorm:"foreignKey:BookingID"`
	SeatID    uuid.UUID `gorm:"type:char(36);not null"`
	Seat      seat.Seat `gorm:"foreignKey:SeatID"`
	Price     int       `gorm:"not null"`
}
