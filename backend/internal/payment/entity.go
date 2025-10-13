package payment

import (
	"github.com/baimhons/stadiumhub/internal/utils"
	"github.com/google/uuid"
)

type Payment struct {
	utils.BaseEntity
	ID uuid.UUID `gorm:"type:char(36);primaryKey"`
	// BookingID uuid.UUID `gorm:"type:char(36);not null"`
	// Booking   booking.Booking `gorm:"foreignKey:BookingID"`
	Amount int    `gorm:"not null"`
	Status string `gorm:"type:varchar(20);not null"`
}
