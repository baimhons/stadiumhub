package booking

import (
	"github.com/baimhons/stadiumhub/internal/utils"
	"gorm.io/gorm"
)

type BookingRepository interface {
	utils.BaseRepository[Booking]
}

type bookingRepositoryImpl struct {
	utils.BaseRepository[Booking]
	db *gorm.DB
}

func NewBookingRepository(db *gorm.DB) BookingRepository {
	return &bookingRepositoryImpl{
		BaseRepository: utils.NewBaseRepository[Booking](db),
		db:             db,
	}
}
