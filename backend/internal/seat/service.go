package seat

import (
	"github.com/google/uuid"
)

type SeatService interface {
	GetAvailableSeats(matchID uint, teamID int, zoneID *uuid.UUID) ([]Seat, error)
}

type seatServiceImpl struct {
	repo SeatRepository
}

func NewSeatService(repo SeatRepository) SeatService {
	return &seatServiceImpl{repo: repo}
}

func (s *seatServiceImpl) GetAvailableSeats(matchID uint, teamID int, zoneID *uuid.UUID) ([]Seat, error) {
	return s.repo.QueryAvailableSeat(matchID, teamID, zoneID)
}
