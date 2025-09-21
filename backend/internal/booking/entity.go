package booking

import (
	"github.com/baimhons/stadiumhub.git/internal/match"
	"github.com/baimhons/stadiumhub.git/internal/user"
	"github.com/baimhons/stadiumhub.git/internal/utils"
	"github.com/baimhons/stadiumhub.git/internal/zone"
	"github.com/google/uuid"
)

type BookingStatus string

const (
	BookingPending  BookingStatus = "PENDING"
	BookingPaid     BookingStatus = "PAID"
	BookingCanceled BookingStatus = "CANCELED"
)

type Booking struct {
	utils.BaseEntity
	UserID     uuid.UUID     `gorm:"not null"`
	User       user.User     `gorm:"foreignKey:UserID"`
	MatchID    uuid.UUID     `gorm:"not null"`
	Match      match.Match   `gorm:"foreignKey:MatchID"`
	ZoneID     uuid.UUID     `gorm:"not null"`
	Zone       zone.Zone     `gorm:"foreignKey:ZoneID"`
	Quantity   int           `gorm:"not null"`
	TotalPrice int           `gorm:"not null"`
	Status     BookingStatus `gorm:"type:varchar(20);not null"`
}
