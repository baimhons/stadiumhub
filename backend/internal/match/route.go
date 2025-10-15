package match

import (
	"github.com/baimhons/stadiumhub/internal/user"
	"github.com/gin-gonic/gin"
)

type MatchRoutes struct {
	group        *gin.RouterGroup
	matchHandler MatchHandler
	userValidate user.UserValidate
}

func NewMatchRoutes(
	group *gin.RouterGroup,
	matchHandler MatchHandler,
	userValidate user.UserValidate,
) *MatchRoutes {

	matchGroup := group.Group("/matches")
	r := &MatchRoutes{
		group:        matchGroup,
		matchHandler: matchHandler,
		userValidate: userValidate,
	}
	return r
}

func (r *MatchRoutes) RegisterRoutes() {
	r.group.GET("/", r.matchHandler.GetAllMatches)
	r.group.GET("/team/:teamID", r.matchHandler.GetMatchesByTeamID)
	r.group.GET("/date-range", r.matchHandler.GetMatchesByDateRange)
	r.group.POST("/update", r.userValidate.ValidateRoleAdmin, r.matchHandler.UpdateMatches)
}
