package zone

import (
	"github.com/baimhons/stadiumhub.git/internal/team"
	"github.com/baimhons/stadiumhub.git/internal/utils"
	"github.com/google/uuid"
)

type Zone struct {
	utils.BaseEntity
	TeamID      uuid.UUID `gorm:"not null"`
	Team        team.Team `gorm:"foreignKey:TeamID"`
	Name        string    `gorm:"not null"`
	Capacity    int       `gorm:"not null"`
	Price       float32   `gorm:"not null"`
	Discription string    `gorm:"not null"`
}
