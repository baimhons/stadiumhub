package seed

import (
	"log"

	"github.com/baimhons/stadiumhub/internal"
	"github.com/baimhons/stadiumhub/internal/user"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// สร้าง admin seed
func SeedAdmin(db *gorm.DB) {
	var existingAdmin user.User
	if err := db.Where("role = ?", "admin").First(&existingAdmin).Error; err == nil {
		log.Println("Admin user already exists")
		return
	} else if err != gorm.ErrRecordNotFound {
		log.Fatalf("failed to check existing admin: %v", err)
	}

	password := internal.ENV.AdminData.Password
	hashPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Fatalf("failed to hash password: %v", err)
	}

	admin := user.User{
		Username: internal.ENV.AdminData.Username,
		Email:    internal.ENV.AdminData.Email,
		Password: string(hashPassword),
		Role:     internal.ENV.AdminData.Role,
	}

	if err := db.Create(&admin).Error; err != nil {
		log.Fatalf("failed to create admin user: %v", err)
	}

	log.Println("Admin user created successfully")
}
