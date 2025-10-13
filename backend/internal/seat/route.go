package seat

import "github.com/gin-gonic/gin"

type SeatRoutes struct {
	group       *gin.RouterGroup
	seatHandler SeatHandler
}

func NewSeatRoutes(
	group *gin.RouterGroup,
	seatHandler SeatHandler,
) *SeatRoutes {

	seatGroup := group.Group("/seat")
	r := &SeatRoutes{
		group:       seatGroup,
		seatHandler: seatHandler,
	}

	return r
}

func (r *SeatRoutes) RegisterRoutes() {

	r.group.GET("/avaliable", r.seatHandler.GetAvailableSeats)
}
