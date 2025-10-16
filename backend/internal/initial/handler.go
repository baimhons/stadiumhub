package initial

import (
	"github.com/baimhons/stadiumhub/internal/booking"
	"github.com/baimhons/stadiumhub/internal/match"
	"github.com/baimhons/stadiumhub/internal/seat"
	"github.com/baimhons/stadiumhub/internal/team"
	"github.com/baimhons/stadiumhub/internal/user"
	"github.com/baimhons/stadiumhub/internal/zone"
)

type Handler struct {
	UserHandler    user.UserHandler
	MatchHandler   match.MatchHandler
	SeatHandler    seat.SeatHandler
	BookingHandler booking.BookingHandler
	ZoneHandler    zone.ZoneHandler
	TeamHandler    team.TeamHandler
}

func NewHandler(service *Service) *Handler {
	return &Handler{
		UserHandler:    user.NewUserHandler(service.UserService),
		MatchHandler:   match.NewMatchHandler(service.MatchService),
		SeatHandler:    seat.NewSeatHandler(service.SeatService),
		BookingHandler: booking.NewBookingHandler(service.BookingService),
		ZoneHandler:    zone.NewZoneHandler(service.ZoneService),
		TeamHandler:    team.NewTeamHandler(service.TeamService),
	}
}
