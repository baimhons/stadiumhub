package initial

import (
	"github.com/baimhons/stadiumhub/internal/match"
	"github.com/baimhons/stadiumhub/internal/user"
)

type Repository struct {
	UserRepository  user.UserRepository
	MatchRepository match.MatchRepository
}

func NewRepository(clientConfig *clientConfig) *Repository {
	return &Repository{
		UserRepository: user.NewUserRepository(
			clientConfig.DB,
		),
		MatchRepository: match.NewMatchRepository(
			clientConfig.DB,
		),
	}
}
