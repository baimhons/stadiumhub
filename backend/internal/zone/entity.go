package zone

import (
	"github.com/baimhons/stadiumhub/internal/team"
	"github.com/baimhons/stadiumhub/internal/utils"
	"github.com/google/uuid"
)

type Zone struct {
	utils.BaseEntity
	TeamID uuid.UUID `gorm:"not null"`
	Team   team.Team `gorm:"foreignKey:TeamID"`
	Name   string    `gorm:"not null"`
}
