package seat

import "github.com/google/uuid"

type SeatService interface {
}
type seatServiceImpl struct {
	seatRepo SeatRepository
}

func NewSeatService(seatRepo SeatRepository) SeatService {
	return &seatServiceImpl{
		seatRepo: seatRepo,
	}
}

func (sr *seatRepositoryImpl) SearchAvaliableSeat(matchID uint, seatID uuid.UUID) ([]Seat, error) {
	return nil, nil
}
