package team

import "github.com/baimhons/stadiumhub.git/internal/utils"

type Team struct {
	utils.BaseEntity
	Name      string `gorm:"not null,unique"`
	ShortName string `gorm:"not null,unique"`
	TLA       string `gorm:"not null"`
	Address   string `gorm:"not null"`
	Venue     string `gorm:"not null,unique"`
}
