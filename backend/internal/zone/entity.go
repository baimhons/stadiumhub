package zone

import (
	"github.com/baimhons/stadiumhub/internal/team"
	"github.com/baimhons/stadiumhub/internal/utils"
)

type Zone struct {
	utils.BaseEntity
	TeamID int       `gorm:"not null"`
	Team   team.Team `gorm:"foreignKey:TeamID"`
	Name   string    `gorm:"not null"`
}
