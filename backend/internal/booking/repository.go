package booking

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/baimhons/stadiumhub/internal/utils"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type BookingRepository interface {
	utils.BaseRepository[Booking]
	GetByIDWithRelations(id uuid.UUID) (*Booking, error)
	GetBookingsByUserID(userID uuid.UUID, query *utils.PaginationQuery) ([]Booking, int, error)
	GetAllWithRelations(pagination *utils.PaginationQuery) ([]Booking, int, error)
	GetRevenueByYear(year int) (map[string]float32, error)
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

func (br *bookingRepositoryImpl) GetByIDWithRelations(id uuid.UUID) (*Booking, error) {
	var booking Booking

	err := br.db.
		Preload("User").
		Preload("Match.HomeTeam").
		Preload("Match.AwayTeam").
		Preload("Seats").
		Preload("Seats.Seat.Zone.Team").
		First(&booking, "id = ?", id).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("booking not found")
		}
		return nil, err
	}

	if len(booking.Seats) == 0 {
		fmt.Println("[DEBUG] booking.Seats empty for booking:", booking.ID)
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
		Preload("Seats.Seat.Zone.Team")

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

func (br *bookingRepositoryImpl) GetRevenueByYear(year int) (map[string]float32, error) {
	revenueMap := make(map[string]float32)
	months := []string{
		"January", "February", "March", "April", "May", "June",
		"July", "August", "September", "October", "November", "December",
	}

	for i, monthName := range months {
		startDate := time.Date(year, time.Month(i+1), 1, 0, 0, 0, 0, time.UTC)
		endDate := startDate.AddDate(0, 1, 0)

		var total *float32
		query := br.db.Model(&Booking{}).
			Where("status = ?", "PAID").
			Where("created_at >= ? AND created_at < ?", startDate, endDate)

		if err := query.Select("SUM(total_price)").Scan(&total).Error; err != nil {
			return nil, err
		}

		if total != nil {
			revenueMap[monthName] = *total
		} else {
			revenueMap[monthName] = 0
		}
	}

	return revenueMap, nil
}
