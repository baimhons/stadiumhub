package initial

import (
	"github.com/baimhons/stadiumhub/internal/booking"
	"github.com/baimhons/stadiumhub/internal/match"
	"github.com/baimhons/stadiumhub/internal/seat"
	"github.com/baimhons/stadiumhub/internal/user"
	"github.com/baimhons/stadiumhub/internal/utils"
)

type Service struct {
	UserService    user.UserService
	MatchService   match.MatchService
	SeatService    seat.SeatService
	BookingService booking.BookingService
}

func NewService(repo *Repository, redis utils.RedisClient) *Service {
	return &Service{
		UserService: user.NewUserService(
			repo.UserRepository,
			redis,
		),
		MatchService: match.NewMatchService(
			repo.MatchRepository,
		),
		SeatService: seat.NewSeatService(
			repo.SeatRepository,
		),
		BookingService: booking.NewBookingService(
			repo.BookingRepository,
		),
	}
}
