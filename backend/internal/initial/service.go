package initial

import (
	"github.com/baimhons/stadiumhub/internal/match"
	"github.com/baimhons/stadiumhub/internal/user"
	"github.com/baimhons/stadiumhub/internal/utils"
)

type Service struct {
	UserService  user.UserService
	MatchService match.MatchService
}

func NewService(repo *Repository, redis utils.RedisClient) *Service {
	return &Service{
		UserService: user.NewUserService(
			repo.UserRepository,
			redis,
		),
		MatchService: match.NewMatchService(
			repo.MatchRepository,
		),
	}
}
