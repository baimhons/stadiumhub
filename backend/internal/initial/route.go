package initial

import (
	"github.com/baimhons/stadiumhub/internal/booking"
	"github.com/baimhons/stadiumhub/internal/match"
	"github.com/baimhons/stadiumhub/internal/seat"
	"github.com/baimhons/stadiumhub/internal/user"
	"github.com/gin-gonic/gin"
)

type route struct {
	UserRoutes    *user.UserRoutes
	MatchRoutes   *match.MatchRoutes
	SeatRoutes    *seat.SeatRoutes
	BookingRoutes *booking.BookingRoutes
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
		),
		SeatRoutes: seat.NewSeatRoutes(
			apiRoute,
			handler.SeatHandler,
		),
		BookingRoutes: booking.NewBookingRoutes(
			apiRoute,
			handler.BookingHandler,
			validate.BookingValidate,
			middleware.AuthMiddleware,
		),
	}
	route.setupRoute()
}

func (r *route) setupRoute() {
	r.UserRoutes.RegisterRoutes()
	r.MatchRoutes.RegisterRoutes()
	r.SeatRoutes.RegisterRoutes()
	r.BookingRoutes.RegisterRoutes()
}
