package match

import (
	"github.com/baimhons/stadiumhub/internal/middlewares"
	"github.com/baimhons/stadiumhub/internal/user"
	"github.com/gin-gonic/gin"
)

type MatchRoutes struct {
	group          *gin.RouterGroup
	matchHandler   MatchHandler
	userValidate   user.UserValidate
	authMiddleware middlewares.AuthMiddlewareImpl
}

func NewMatchRoutes(
	group *gin.RouterGroup,
	matchHandler MatchHandler,
	userValidate user.UserValidate,
	authMiddleware middlewares.AuthMiddlewareImpl,
) *MatchRoutes {

	matchGroup := group.Group("/matches")
	r := &MatchRoutes{
		group:          matchGroup,
		matchHandler:   matchHandler,
		userValidate:   userValidate,
		authMiddleware: authMiddleware,
	}
	return r
}

func (r *MatchRoutes) RegisterRoutes() {
	r.group.GET("/", r.matchHandler.GetAllMatches)
	r.group.GET("/team/:teamID", r.matchHandler.GetMatchesByTeamID)
	r.group.GET("/date-range", r.matchHandler.GetMatchesByDateRange)
	r.group.POST("/update", r.authMiddleware.RequireAuth(), r.userValidate.ValidateRoleAdmin, r.matchHandler.UpdateMatches)
}
