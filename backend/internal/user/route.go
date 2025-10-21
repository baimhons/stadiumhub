package user

import (
	"github.com/baimhons/stadiumhub/internal/middlewares"
	"github.com/gin-gonic/gin"
)

type UserRoutes struct {
	group          *gin.RouterGroup
	userHandler    UserHandler
	userValidate   UserValidate
	authMiddleware middlewares.AuthMiddlewareImpl
}

func NewUserRoutes(
	group *gin.RouterGroup,
	userHandler UserHandler,
	userValidate UserValidate,
	authMiddleware middlewares.AuthMiddlewareImpl,
) *UserRoutes {

	userGroup := group.Group("/user")
	r := &UserRoutes{
		group:          userGroup,
		userHandler:    userHandler,
		userValidate:   userValidate,
		authMiddleware: authMiddleware,
	}

	return r
}

func (r *UserRoutes) RegisterRoutes() {

	r.group.GET("/profile", r.authMiddleware.RequireAuth(), r.userHandler.GetUserProfile)
	r.group.POST("/register", r.userValidate.ValidateRegisterUser, r.userHandler.RegisterUser)
	r.group.POST("/login", r.userValidate.ValidateLoginUser, r.userHandler.LoginUser)
	r.group.POST("/logout", r.authMiddleware.RequireAuth(), r.userHandler.LogoutUser)
	r.group.PUT("/update", r.authMiddleware.RequireAuth(), r.userValidate.ValidateUpdateUser, r.userHandler.UpdateUser)
}
