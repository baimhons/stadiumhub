package initial

import (
	"github.com/baimhons/stadiumhub/internal/user"
	"github.com/gin-gonic/gin"
)

type route struct {
	UserRoutes *user.UserRoutes
}

func NewRoute(
	engine *gin.Engine,
	handler Handler,
	validate Validate,
) {
	apiRoute := engine.Group("/api/v1")
	route := &route{
		UserRoutes: user.NewUserRoutes(
			apiRoute,
			handler.UserHandler,
			validate.UserValidate,
		),
	}
	route.setupRoute()
}

func (r *route) setupRoute() {
	r.UserRoutes.RegisterRoutes()
}
