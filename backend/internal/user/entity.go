package user

import (
	"github.com/baimhons/stadiumhub/internal/utils"
)

type UserRole string

// const (
// 	RoleAdmin UserRole = "ADMIN"
// 	RoleUser  UserRole = "USER"
// )

type User struct {
	utils.BaseEntity
	Username    string `gorm:"not null;unique"`
	FirstName   string `gorm:"not null"`
	LastName    string `gorm:"not null"`
	Email       string `gorm:"not null;unique"`
	Password    string `gorm:"not null"`
	PhoneNumber string `gorm:"not null;unique"`
	Role        string `gorm:"type:varchar(20);not null"`
}
