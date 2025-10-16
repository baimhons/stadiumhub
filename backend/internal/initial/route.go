package initial

import (
	"github.com/baimhons/stadiumhub/internal/booking"
	"github.com/baimhons/stadiumhub/internal/match"
	"github.com/baimhons/stadiumhub/internal/seat"
	"github.com/baimhons/stadiumhub/internal/team"
	"github.com/baimhons/stadiumhub/internal/user"
	"github.com/baimhons/stadiumhub/internal/zone"
	"github.com/gin-gonic/gin"
)

type route struct {
	UserRoutes    *user.UserRoutes
	MatchRoutes   *match.MatchRoutes
	SeatRoutes    *seat.SeatRoutes
	BookingRoutes *booking.BookingRoutes
	ZoneRoutes    *zone.ZoneRoutes
	TeamRoutes    *team.TeamRoutes
}

func NewRoute(
	engine *gin.Engine,
	handler Handler,
	validate Validate,
	middleware Middleware,
) {
	apiRoute := engine.Group("/api/v1")
	route := &route{
		UserRoutes: user.NewUserRoutes(
			apiRoute,
			handler.UserHandler,
			validate.UserValidate,
			middleware.AuthMiddleware,
		),
		MatchRoutes: match.NewMatchRoutes(
			apiRoute,
			handler.MatchHandler,
			validate.UserValidate,
			middleware.AuthMiddleware,
		),
		SeatRoutes: seat.NewSeatRoutes(
			apiRoute,
			handler.SeatHandler,
		),
		BookingRoutes: booking.NewBookingRoutes(
			apiRoute,
			handler.BookingHandler,
			validate.BookingValidate,
			validate.UserValidate,
			middleware.AuthMiddleware,
		),
		ZoneRoutes: zone.NewZoneRoutes(
			apiRoute,
			handler.ZoneHandler,
		),
		TeamRoutes: team.NewTeamRoutes(
			apiRoute,
			handler.TeamHandler,
		),
	}
	route.setupRoute()
}

func (r *route) setupRoute() {
	r.UserRoutes.RegisterRoutes()
	r.MatchRoutes.RegisterRoutes()
	r.SeatRoutes.RegisterRoutes()
	r.BookingRoutes.RegisterRoutes()
	r.ZoneRoutes.RegisterRoutes()
	r.TeamRoutes.RegisterRoutes()
}
