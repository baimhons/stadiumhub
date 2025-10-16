package team

import "github.com/gin-gonic/gin"

type TeamRoutes struct {
	group       *gin.RouterGroup
	teamHandler TeamHandler
}

func NewTeamRoutes(
	group *gin.RouterGroup,
	teamHandler TeamHandler,
) *TeamRoutes {

	teamGroup := group.Group("/team")
	r := &TeamRoutes{
		group:       teamGroup,
		teamHandler: teamHandler,
	}

	return r
}

func (r *TeamRoutes) RegisterRoutes() {

	r.group.GET("/", r.teamHandler.GetAllTeam)
}
