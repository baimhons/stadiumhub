package zone

import "github.com/gin-gonic/gin"

type ZoneRoutes struct {
	group       *gin.RouterGroup
	zoneHandler ZoneHandler
}

func NewZoneRoutes(
	group *gin.RouterGroup,
	zoneHandler ZoneHandler,
) *ZoneRoutes {

	zoneGroup := group.Group("/zone")
	r := &ZoneRoutes{
		group:       zoneGroup,
		zoneHandler: zoneHandler,
	}

	return r
}

func (r *ZoneRoutes) RegisterRoutes() {
	r.group.GET("/:id", r.zoneHandler.GetZoneByMatchID)
}
