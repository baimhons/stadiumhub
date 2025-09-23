package initial

import (
	"github.com/baimhons/stadiumhub/internal/user"
)

type Repository struct {
	UserRepository user.UserRepository
}

func NewRepository(clientConfig *clientConfig) *Repository {
	return &Repository{
		UserRepository: user.NewUserRepository(
			clientConfig.DB,
		),
	}
}
