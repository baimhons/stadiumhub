package initial

import (
	"github.com/baimhons/stadiumhub/internal/match"
	"github.com/baimhons/stadiumhub/internal/user"
)

type Handler struct {
	UserHandler  user.UserHandler
	MatchHandler match.MatchHandler
}

func NewHandler(service *Service) *Handler {
	return &Handler{
		UserHandler:  user.NewUserHandler(service.UserService),
		MatchHandler: match.NewMatchHandler(service.MatchService),
	}
}
