package initial

import (
	"github.com/baimhons/stadiumhub/internal/booking"
	"github.com/baimhons/stadiumhub/internal/match"
	"github.com/baimhons/stadiumhub/internal/seat"
	"github.com/baimhons/stadiumhub/internal/user"
)

type Handler struct {
	UserHandler    user.UserHandler
	MatchHandler   match.MatchHandler
	SeatHandler    seat.SeatHandler
	BookingHandler booking.BookingHandler
}

func NewHandler(service *Service) *Handler {
	return &Handler{
		UserHandler:    user.NewUserHandler(service.UserService),
		MatchHandler:   match.NewMatchHandler(service.MatchService),
		SeatHandler:    seat.NewSeatHandler(service.SeatService),
		BookingHandler: booking.NewBookingHandler(service.BookingService),
	}
}
