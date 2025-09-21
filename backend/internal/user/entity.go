package user

import "github.com/baimhons/stadiumhub.git/internal/utils"

type UserRole string

const (
	RoleAdmin UserRole = "ADMIN"
	RoleUser  UserRole = "USER"
)

type User struct {
	utils.BaseEntity
	Username    string   `gorm:"not null;unique"`
	Firstname   string   `gorm:"not null"`
	Lastname    string   `gorm:"not null"`
	Email       string   `gorm:"not null;unique"`
	Password    string   `gorm:"not null"`
	PhoneNumber string   `gorm:"not null;unique"`
	Role        UserRole `gorm:"type:varchar(20);not null"`
}
