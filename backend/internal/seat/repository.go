package seat

import (
	"github.com/baimhons/stadiumhub/internal/utils"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type SeatRepository interface {
	utils.BaseRepository[Seat]
	QueryAvailableSeat(matchID uint, teamID int, zoneID *uuid.UUID) ([]Seat, error)
}

type seatRepositoryImpl struct {
	utils.BaseRepository[Seat]
	db *gorm.DB
}

func NewSeatRepository(db *gorm.DB) SeatRepository {
	return &seatRepositoryImpl{
		BaseRepository: utils.NewBaseRepository[Seat](db),
		db:             db,
	}
}

func (sr *seatRepositoryImpl) QueryAvailableSeat(matchID uint, teamID int, zoneID *uuid.UUID) ([]Seat, error) {
	var seats []Seat

	query := sr.db.
		Model(&Seat{}).
		Joins("JOIN zones z ON seats.zone_id = z.id").
		Where("z.team_id = ?", teamID).
		Where(`seats.id NOT IN (
			SELECT bs.seat_id
			FROM booking_seats bs
			JOIN bookings b ON bs.booking_id = b.id
			WHERE b.match_id = ?
			AND b.status IN ('PENDING', 'PAID')
		)`, matchID)

	if zoneID != nil {
		query = query.Where("seats.zone_id = ?", *zoneID)
	}

	if err := query.Find(&seats).Error; err != nil {
		return nil, err
	}

	return seats, nil
}
