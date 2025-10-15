package booking

import (
	"errors"
	"net/http"

	"github.com/baimhons/stadiumhub/internal/booking/api/request"
	"github.com/baimhons/stadiumhub/internal/booking/api/response"
	"github.com/baimhons/stadiumhub/internal/match"
	"github.com/baimhons/stadiumhub/internal/models"
	"github.com/baimhons/stadiumhub/internal/seat"
	"github.com/google/uuid"
)

type BookingService interface {
	CreateBooking(userCtx models.UserContext, req request.CreateBookingRequest) (resp *response.BookingResponse, statusCode int, err error)
	GetBookingByID(id uuid.UUID) (resp *response.BookingResponse, statusCode int, err error)
	GetAllBookingsByUser(userID uuid.UUID) (resp []response.BookingResponse, statusCode int, err error)
	CancelBooking(id uuid.UUID) (statusCode int, err error)
	GetAllBookings() (resp []response.BookingResponse, statusCode int, err error)
	UpdateBookingStatus(req request.UpdateBookingStatusRequest) (resp *response.BookingResponse, statusCode int, err error)
}
type bookingServiceImpl struct {
	bookingRepository BookingRepository
}

func NewBookingService(bookingRepository BookingRepository) BookingService {
	return &bookingServiceImpl{
		bookingRepository: bookingRepository,
	}
}

// create booking service
func (bs *bookingServiceImpl) CreateBooking(userCtx models.UserContext, req request.CreateBookingRequest) (resp *response.BookingResponse, statusCode int, err error) {
	if userCtx.ID == uuid.Nil {
		return nil, http.StatusBadRequest, errors.New("user context is required")
	}

	tx := bs.bookingRepository.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()
	if tx.Error != nil {
		return nil, http.StatusBadRequest, tx.Error
	}
	var match match.Match
	if err := tx.Preload("HomeTeam").Preload("AwayTeam").First(&match, req.MatchID).Error; err != nil {
		tx.Rollback()
		return nil, http.StatusNotFound, err
	}

	seatPrice := match.HomeTeam.Price

	var validSeats []seat.Seat
	if err := tx.Joins("JOIN zones z ON seats.zone_id = z.id").
		Where("seats.id IN ?", req.SeatIDs).
		Where("z.team_id = ?", match.HomeTeam.ID).
		Find(&validSeats).Error; err != nil {
		tx.Rollback()
		return nil, http.StatusInternalServerError, err
	}

	if len(validSeats) != len(req.SeatIDs) {
		tx.Rollback()
		return nil, http.StatusBadRequest, errors.New("some selected seats are not valid for this match")
	}

	var bookedSeats []BookingSeat
	if err := tx.Where("seat_id IN ? AND booking_id IN (SELECT id FROM bookings WHERE match_id = ? AND status != ?)", req.SeatIDs, req.MatchID, "CANCELED").
		Find(&bookedSeats).Error; err != nil {
		tx.Rollback()
		return nil, http.StatusInternalServerError, err
	}

	if len(bookedSeats) > 0 {
		tx.Rollback()
		return nil, http.StatusBadRequest, errors.New("some selected seats are already booked")
	}

	totalPrice := float32(len(req.SeatIDs)) * seatPrice
	newBooking := Booking{
		UserID:     userCtx.ID,
		MatchID:    req.MatchID,
		TotalPrice: int(totalPrice),
		Status:     "PENDING",
	}

	if err := tx.Create(&newBooking).Error; err != nil {
		tx.Rollback()
		return nil, http.StatusBadRequest, err
	}

	var bookingSeats []BookingSeat
	for _, s := range validSeats {
		bookingSeats = append(bookingSeats, BookingSeat{
			BookingID: newBooking.ID,
			SeatID:    s.ID,
			Price:     int(seatPrice),
		})
	}

	if err := tx.Create(&bookingSeats).Error; err != nil {
		tx.Rollback()
		return nil, http.StatusBadRequest, err
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return nil, http.StatusBadRequest, err
	}

	var seatResp []response.BookingSeatResp
	for _, bs := range bookingSeats {
		seatResp = append(seatResp, response.BookingSeatResp{
			SeatID: bs.SeatID,
			Price:  bs.Price,
		})
	}

	resp = &response.BookingResponse{
		ID:         newBooking.ID,
		UserID:     newBooking.UserID,
		MatchID:    newBooking.MatchID,
		TotalPrice: newBooking.TotalPrice,
		Status:     newBooking.Status,
		Seats:      seatResp,
		CreatedAt:  newBooking.CreatedAt,
		UpdatedAt:  newBooking.UpdatedAt,
	}

	return resp, http.StatusCreated, nil
}

// get booking by id
func (bs *bookingServiceImpl) GetBookingByID(id uuid.UUID) (resp *response.BookingResponse, statusCode int, err error) {
	return nil, 0, nil
}

// get all bookings (user)
func (bs *bookingServiceImpl) GetAllBookingsByUser(userID uuid.UUID) (resp []response.BookingResponse, statusCode int, err error) {
	return nil, 0, nil
}

// cancel booking
func (bs *bookingServiceImpl) CancelBooking(id uuid.UUID) (statusCode int, err error) {
	return 0, nil
}

// get all bookings (admin)
func (bs *bookingServiceImpl) GetAllBookings() (resp []response.BookingResponse, statusCode int, err error) {
	return nil, 0, nil
}

// update booking status
func (bs *bookingServiceImpl) UpdateBookingStatus(req request.UpdateBookingStatusRequest) (resp *response.BookingResponse, statusCode int, err error) {
	return nil, 0, nil
}
