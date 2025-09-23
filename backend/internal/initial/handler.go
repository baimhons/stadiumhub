package initial

import "github.com/baimhons/stadiumhub/internal/user"

type Handler struct {
	UserHandler user.UserHandler
}

func NewHandler(service *Service) *Handler {
	return &Handler{
		UserHandler: user.NewUserHandler(service.UserService),
	}
}
