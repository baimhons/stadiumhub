package user

import (
	"github.com/baimhons/stadiumhub/internal/utils"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type UserRole string

// const (
// 	RoleAdmin UserRole = "ADMIN"
// 	RoleUser  UserRole = "USER"
// )

type User struct {
	utils.BaseEntity
	Username    string `gorm:"not null;unique"`
	FullName    string `gorm:"not null"`
	Email       string `gorm:"not null;unique"`
	Password    string `gorm:"not null"`
	PhoneNumber string `gorm:"not null;unique"`
	Role        string `gorm:"type:varchar(20);not null"`
}

type TokenContext struct {
	UserID uuid.UUID `json:"user_id"`
	jwt.RegisteredClaims
}

type UserContext struct {
	UserID       uuid.UUID `json:"user_id"`
	Username     string    `json:"username"`
	Email        string    `json:"email"`
	FullName     string    `json:"full_name"`
	Role         string    `json:"role"`
	AccessToken  string    `json:"access_token"`
	RefreshToken string    `json:"refresh_token"`
}
