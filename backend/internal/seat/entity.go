package seat

import (
	"github.com/baimhons/stadiumhub/internal/utils"
	"github.com/baimhons/stadiumhub/internal/zone"
	"github.com/google/uuid"
)

// type SeatStatus string

// const (
// 	SeatAvailable SeatStatus = "AVAILABLE"
// 	SeatReserved  SeatStatus = "RESERVED"
// )

type Seat struct {
	utils.BaseEntity
	SeatNo string    `gorm:"not null,unique"`
	ZoneID uuid.UUID `gorm:"type:char(36);not null"`
	Zone   zone.Zone `gorm:"foreignKey:ZoneID"`
}
