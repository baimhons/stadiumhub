package user

import (
	"github.com/gin-gonic/gin"
)

type UserRoutes struct {
	group        *gin.RouterGroup
	userHandler  UserHandler
	userValidate UserValidate
}

func NewUserRoutes(
	group *gin.RouterGroup,
	userHandler UserHandler,
	userValidate UserValidate,
) *UserRoutes {
	r := &UserRoutes{
		group:        group,
		userHandler:  userHandler,
		userValidate: userValidate,
	}

	r.RegisterRoutes()

	return r
}

func (r *UserRoutes) RegisterRoutes() {
	r.group.POST("/register", r.userValidate.ValidateRegisterUser, r.userHandler.RegisterUser)
	r.group.POST("/login", r.userValidate.ValidateLoginUser, r.userHandler.LoginUser)
}
