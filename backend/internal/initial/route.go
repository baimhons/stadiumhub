package initial

import (
	"github.com/baimhons/stadiumhub/internal/match"
	"github.com/baimhons/stadiumhub/internal/user"
	"github.com/gin-gonic/gin"
)

type route struct {
	UserRoutes  *user.UserRoutes
	MatchRoutes *match.MatchRoutes
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
		MatchRoutes: match.NewMatchRoutes(
			apiRoute,
			handler.MatchHandler,
		),
	}
	route.setupRoute()
}

func (r *route) setupRoute() {
	r.UserRoutes.RegisterRoutes()
	r.MatchRoutes.RegisterRoutes()
}
