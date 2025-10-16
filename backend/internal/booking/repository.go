package booking

import (
	"fmt"
	"net/http"

	"github.com/baimhons/stadiumhub/internal/utils"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type BookingRepository interface {
	utils.BaseRepository[Booking]
	GetByIDWithRelations(id uuid.UUID) (*Booking, error)
	GetBookingsByUserID(userID uuid.UUID, query *utils.PaginationQuery) ([]Booking, int, error)
	GetAllWithRelations(pagination *utils.PaginationQuery) ([]Booking, int, error)
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

func (br *bookingRepositoryImpl) GetBookingsByUserID(userID uuid.UUID, query *utils.PaginationQuery) ([]Booking, int, error) {
	var bookings []Booking

	tx := br.db.
		Model(&Booking{}).
		Where("user_id = ?", userID).
		Preload("User").
		Preload("Seats").
		Preload("Match").
		Preload("Match.HomeTeam").
		Preload("Match.AwayTeam")

	if query.Page != nil && query.PageSize != nil {
		offset := (*query.Page - 1) * (*query.PageSize)
		tx = tx.Offset(offset).Limit(*query.PageSize)
	}

	if query.Sort != nil && query.Order != nil {
		tx = tx.Order(*query.Sort + " " + *query.Order)
	}

	if err := tx.Find(&bookings).Error; err != nil {
		return nil, http.StatusInternalServerError, err
	}

	if len(bookings) == 0 {
		return nil, http.StatusNotFound, nil
	}

	return bookings, http.StatusOK, nil
}

func (r *bookingRepositoryImpl) GetByIDWithRelations(id uuid.UUID) (*Booking, error) {
	var booking Booking
	if err := r.db.
		Preload("User").
		Preload("Match").
		Preload("Match.HomeTeam").
		Preload("Match.AwayTeam").
		Preload("Seats").
		First(&booking, "id = ?", id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("booking not found")
		}
		return nil, err
	}

	return &booking, nil
}

func (br *bookingRepositoryImpl) GetAllWithRelations(pagination *utils.PaginationQuery) ([]Booking, int, error) {
	var bookings []Booking

	query := br.db.
		Model(&Booking{}).
		Preload("User").
		Preload("Match.HomeTeam").
		Preload("Match.AwayTeam").
		Preload("BookingSeats.Seat.Zone.Team")

	if pagination.Page != nil && pagination.PageSize != nil {
		offset := (*pagination.Page - 1) * (*pagination.PageSize)
		query = query.Offset(offset).Limit(*pagination.PageSize)
	}

	if pagination.Sort != nil && pagination.Order != nil {
		query = query.Order(*pagination.Sort + " " + *pagination.Order)
	}

	if err := query.Find(&bookings).Error; err != nil {
		return nil, http.StatusInternalServerError, err
	}

	if len(bookings) == 0 {
		return nil, http.StatusNotFound, nil
	}

	return bookings, http.StatusOK, nil
}
