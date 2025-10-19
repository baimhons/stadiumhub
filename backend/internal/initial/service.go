package initial

import (
	"github.com/baimhons/stadiumhub/internal/booking"
	"github.com/baimhons/stadiumhub/internal/match"
	"github.com/baimhons/stadiumhub/internal/payment"
	"github.com/baimhons/stadiumhub/internal/seat"
	"github.com/baimhons/stadiumhub/internal/team"
	"github.com/baimhons/stadiumhub/internal/user"
	"github.com/baimhons/stadiumhub/internal/utils"
	"github.com/baimhons/stadiumhub/internal/zone"
)

type Service struct {
	UserService    user.UserService
	MatchService   match.MatchService
	SeatService    seat.SeatService
	BookingService booking.BookingService
	ZoneService    zone.ZoneService
	TeamService    team.TeamService
	PaymentService payment.PaymentService
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
		ZoneService: zone.NewZoneService(
			repo.ZoneRepository,
		),
		TeamService: team.NewTeamService(
			repo.TeamRepository,
		),
		PaymentService: payment.NewPaymentService(
			repo.BookingRepository,
		),
	}
}
