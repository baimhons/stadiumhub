package user

import (
	"github.com/baimhons/stadiumhub/internal/utils"
	"gorm.io/gorm"
)

type UserRepository interface {
	utils.BaseRepository[User]
	GetByUsernameOrEmail(usernameOrEmail string, item *User) error
}

type userRepositoryImpl struct {
	utils.BaseRepository[User]
	DB *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepositoryImpl{
		BaseRepository: utils.NewBaseRepository[User](db),
		DB:             db,
	}
}

func (r *userRepositoryImpl) GetByUsernameOrEmail(usernameOrEmail string, item *User) error {
	return r.DB.Where("username = ? OR email = ?", usernameOrEmail, usernameOrEmail).First(item).Error
}
