package seat

import (
	"github.com/baimhons/stadiumhub/internal/utils"
)

type SeatRepository interface {
	utils.BaseRepository[Seat]
}

type seatRepositoryImpl struct {
	utils.BaseRepository[Seat]
}

func NewSeatRepository() SeatRepository {
	return &seatRepositoryImpl{
		BaseRepository: utils.NewBaseRepository[Seat](nil),
	}
}
