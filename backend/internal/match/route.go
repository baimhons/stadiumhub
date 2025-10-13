package match

import "github.com/gin-gonic/gin"

type MatchRoutes struct {
	group        *gin.RouterGroup
	MatchHandler MatchHandler
}

func NewMatchRoutes(
	group *gin.RouterGroup,
	matchHandler MatchHandler,
) *MatchRoutes {

	matchGroup := group.Group("/matches")
	r := &MatchRoutes{
		group:        matchGroup,
		MatchHandler: matchHandler,
	}
	return r
}

func (r *MatchRoutes) RegisterRoutes() {
	r.group.GET("/", r.MatchHandler.GetAllMatches)
	r.group.GET("/team/:teamID", r.MatchHandler.GetMatchesByTeamID)
	r.group.GET("/date-range", r.MatchHandler.GetMatchesByDateRange)
	r.group.POST("/update", r.MatchHandler.UpdateMatches)
}
