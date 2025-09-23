package initial

import (
	"github.com/baimhons/stadiumhub/internal/user"
	"github.com/baimhons/stadiumhub/internal/utils"
)

type Service struct {
	UserService user.UserService
}

func NewService(repo *Repository, redis utils.RedisClient) *Service {
	return &Service{
		UserService: user.NewUserService(
			repo.UserRepository,
			redis,
		),
	}
}
