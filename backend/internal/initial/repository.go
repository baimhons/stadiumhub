package initial

import (
	"github.com/baimhons/stadiumhub/internal/booking"
	"github.com/baimhons/stadiumhub/internal/match"
	"github.com/baimhons/stadiumhub/internal/seat"
	"github.com/baimhons/stadiumhub/internal/user"
)

type Repository struct {
	UserRepository    user.UserRepository
	MatchRepository   match.MatchRepository
	SeatRepository    seat.SeatRepository
	BookingRepository booking.BookingRepository
}

func NewRepository(clientConfig *clientConfig) *Repository {
	return &Repository{
		UserRepository: user.NewUserRepository(
			clientConfig.DB,
		),
		MatchRepository: match.NewMatchRepository(
			clientConfig.DB,
		),
		SeatRepository: seat.NewSeatRepository(
			clientConfig.DB,
		),
		BookingRepository: booking.NewBookingRepository(
			clientConfig.DB,
		),
	}
}
