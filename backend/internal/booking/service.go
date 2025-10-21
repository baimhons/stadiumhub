package booking

import (
	"errors"
	"fmt"
	"net/http"
	"net/smtp"
	"strings"
	"time"

	"github.com/baimhons/stadiumhub/internal"
	"github.com/baimhons/stadiumhub/internal/booking/api/request"
	"github.com/baimhons/stadiumhub/internal/booking/api/response"
	"github.com/baimhons/stadiumhub/internal/match"
	"github.com/baimhons/stadiumhub/internal/models"
	"github.com/baimhons/stadiumhub/internal/seat"
	"github.com/baimhons/stadiumhub/internal/utils"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type BookingService interface {
	CreateBooking(userCtx models.UserContext, req request.CreateBookingRequest) (resp *response.BookingResponse, statusCode int, err error)
	GetBookingByID(id uuid.UUID, userCtx models.UserContext) (resp *response.BookingResponse, statusCode int, err error)
	GetAllBookingsByUser(userID uuid.UUID, query *utils.PaginationQuery) (resp []response.BookingResponse, statusCode int, err error)
	CancelBooking(userID uuid.UUID, id uuid.UUID) (statusCode int, err error)
	GetAllBookings(query *utils.PaginationQuery) (resp []response.BookingResponse, statusCode int, err error)
	UpdateBookingStatus(userID uuid.UUID, id uuid.UUID) (statusCode int, err error)
	GetRevenueByYear(year int) (resp []response.MonthRevenue, err error)
	CancelExpiredBookings(expireDuration time.Duration) (rows int, err error)
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
		TotalPrice: totalPrice,
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
			SeatNo:    s.SeatNo,
			Price:     seatPrice,
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
	for _, s := range validSeats {
		seatResp = append(seatResp, response.BookingSeatResp{
			SeatID: s.ID,
			SeatNo: s.SeatNo,
			Price:  seatPrice,
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

	go bs.sendBookingConfirmationEmail(userCtx.Email, newBooking, match, seatResp)

	return resp, http.StatusCreated, nil
}

func (bs *bookingServiceImpl) sendBookingConfirmationEmail(userEmail string, booking Booking, match match.Match, seats []response.BookingSeatResp) {
	from := "stadiumhubtest@gmail.com"
	to := []string{userEmail}
	appPassword := internal.ENV.EmailKey.EmailKey
	host := "smtp.gmail.com"
	addr := host + ":587"

	auth := smtp.PlainAuth("", from, strings.ReplaceAll(appPassword, " ", ""), host)

	seatList := ""
	for _, s := range seats {
		seatList += fmt.Sprintf("  - Seat No: %s (Price: $%.2f)\n", s.SeatNo, s.Price)
	}

	expiryTime := booking.CreatedAt.Add(30 * time.Minute)

	msg := ""
	msg += fmt.Sprintf("From: Stadium Hub <%s>\r\n", from)
	msg += fmt.Sprintf("To: %s\r\n", strings.Join(to, ","))
	msg += "Subject: Booking Confirmation - Payment Required\r\n"
	msg += "MIME-Version: 1.0\r\n"
	msg += "Content-Type: text/plain; charset=UTF-8\r\n"
	msg += "Content-Transfer-Encoding: 8bit\r\n\r\n"

	msg += "Dear Customer,\n\n"
	msg += "Your booking has been successfully created!\n\n"
	msg += "=== BOOKING DETAILS ===\n"
	msg += fmt.Sprintf("Booking ID: %s\n", booking.ID)
	msg += fmt.Sprintf("Match: %s vs %s\n", match.HomeTeam.Name, match.AwayTeam.Name)
	msg += fmt.Sprintf("Match Date: %s\n", match.UTCDate.Format("January 02, 2006 15:04 MST"))
	msg += fmt.Sprintf("Total Price: $%.2f\n\n", booking.TotalPrice)

	msg += "=== YOUR SEATS ===\n"
	msg += seatList
	msg += "\n"

	msg += "⚠️ IMPORTANT: PAYMENT REQUIRED ⚠️\n"
	msg += fmt.Sprintf("Please complete your payment within 30 minutes (before %s)\n", expiryTime.Format("January 02, 2006 15:04 MST"))
	msg += "If payment is not received within this time, your booking will be automatically cancelled.\n\n"

	msg += "To complete your payment, please visit:\n"
	msg += "[Payment URL - Add your payment link here]\n\n"

	msg += "Thank you for choosing Stadium Hub!\n\n"
	msg += "Best regards,\n"
	msg += "Stadium Hub Team"

	if err := smtp.SendMail(addr, auth, from, to, []byte(msg)); err != nil {
		fmt.Printf("Failed to send confirmation email: %v\n", err)
	}
}

// get booking by id
func (bs *bookingServiceImpl) GetBookingByID(id uuid.UUID, userCtx models.UserContext) (resp *response.BookingResponse, statusCode int, err error) {
	booking, err := bs.bookingRepository.GetByIDWithRelations(id)
	if err != nil {
		if err.Error() == "booking not found" {
			return nil, http.StatusNotFound, err
		}
		return nil, http.StatusInternalServerError, err
	}
	if booking.UserID != userCtx.ID {
		return nil, http.StatusForbidden, fmt.Errorf("you are not allowed to access this booking")
	}

	var seatResp []response.BookingSeatResp
	for _, bs := range booking.Seats {
		seatResp = append(seatResp, response.BookingSeatResp{
			SeatID: bs.SeatID,
			SeatNo: bs.SeatNo,
			Price:  bs.Price,
		})
	}

	booking.User.Password = ""

	resp = &response.BookingResponse{
		ID:         booking.ID,
		UserID:     booking.UserID,
		User:       booking.User,
		MatchID:    booking.MatchID,
		Match:      booking.Match,
		TotalPrice: booking.TotalPrice,
		Status:     booking.Status,
		Seats:      seatResp,
		CreatedAt:  booking.CreatedAt,
		UpdatedAt:  booking.UpdatedAt,
	}

	return resp, http.StatusOK, nil
}

// get all bookings (user)
func (bs *bookingServiceImpl) GetAllBookingsByUser(userID uuid.UUID, query *utils.PaginationQuery) (resp []response.BookingResponse, statusCode int, err error) {
	bookings, statusCode, err := bs.bookingRepository.GetBookingsByUserID(userID, query)
	if err != nil {
		return nil, statusCode, err
	}
	if len(bookings) == 0 {
		return nil, http.StatusNotFound, errors.New("no bookings found")
	}

	for _, b := range bookings {

		b.User.Password = ""
		resp = append(resp, response.BookingResponse{
			ID:         b.ID,
			UserID:     b.UserID,
			User:       b.User,
			MatchID:    b.MatchID,
			Match:      b.Match,
			TotalPrice: b.TotalPrice,
			Status:     b.Status,
			CreatedAt:  b.CreatedAt,
			UpdatedAt:  b.UpdatedAt,
		})
	}

	return resp, http.StatusOK, nil
}

func (bs *bookingServiceImpl) CancelBooking(userID uuid.UUID, id uuid.UUID) (statusCode int, err error) {

	booking, err := bs.bookingRepository.GetByIDWithRelations(id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return http.StatusNotFound, errors.New("booking not found")
		}
		return http.StatusInternalServerError, err
	}

	if userID != booking.User.ID {
		// fmt.Println("userId : ", userID)
		// fmt.Println("booking user id : ", booking.User.ID)
		return http.StatusBadRequest, errors.New("booking cannot cancel by another user")
	}

	if booking.Status == "CANCELED" {
		return http.StatusBadRequest, errors.New("booking already cancelled")
	}

	booking.Status = "CANCELED"

	if err := bs.bookingRepository.Update(booking); err != nil {
		return http.StatusInternalServerError, err
	}

	go bs.sendBookingCancellationEmail(booking)

	return http.StatusOK, nil
}

func (bs *bookingServiceImpl) sendBookingCancellationEmail(booking *Booking) {
	from := "stadiumhubtest@gmail.com"
	to := []string{booking.User.Email}
	appPassword := internal.ENV.EmailKey.EmailKey
	host := "smtp.gmail.com"
	addr := host + ":587"

	auth := smtp.PlainAuth("", from, strings.ReplaceAll(appPassword, " ", ""), host)

	// สร้างรายการที่นั่ง
	seatList := ""
	for _, s := range booking.Seats {
		seatList += fmt.Sprintf("  - Seat No: %s (Price: $%.2f)\n", s.SeatNo, s.Price)
	}

	// สร้าง email message
	msg := ""
	msg += fmt.Sprintf("From: Stadium Hub <%s>\r\n", from)
	msg += fmt.Sprintf("To: %s\r\n", strings.Join(to, ","))
	msg += "Subject: Booking Cancellation Confirmation\r\n"
	msg += "MIME-Version: 1.0\r\n"
	msg += "Content-Type: text/plain; charset=UTF-8\r\n"
	msg += "Content-Transfer-Encoding: 8bit\r\n\r\n"

	msg += "Dear Customer,\n\n"
	msg += "Your booking has been successfully cancelled.\n\n"
	msg += "=== CANCELLED BOOKING DETAILS ===\n"
	msg += fmt.Sprintf("Booking ID: %s\n", booking.ID)
	msg += fmt.Sprintf("Match: %s vs %s\n", booking.Match.HomeTeam.Name, booking.Match.AwayTeam.Name)
	msg += fmt.Sprintf("Match Date: %s\n", booking.Match.UTCDate.Format("January 02, 2006 15:04 MST"))
	msg += fmt.Sprintf("Total Amount: $%.2f\n\n", booking.TotalPrice)

	msg += "=== CANCELLED SEATS ===\n"
	msg += seatList
	msg += "\n"

	if booking.Status == "PAID" {
		msg += "Your refund will be processed within 7-14 business days.\n\n"
	}

	msg += "If you did not request this cancellation, please contact us immediately.\n\n"
	msg += "Thank you for using Stadium Hub.\n\n"
	msg += "Best regards,\n"
	msg += "Stadium Hub Team"

	// ส่งอีเมล
	if err := smtp.SendMail(addr, auth, from, to, []byte(msg)); err != nil {
		fmt.Printf("Failed to send cancellation email: %v\n", err)
	}
}

func (bs *bookingServiceImpl) GetAllBookings(query *utils.PaginationQuery) (resp []response.BookingResponse, statusCode int, err error) {
	bookings, statusCode, err := bs.bookingRepository.GetAllWithRelations(query)
	if err != nil {
		return nil, statusCode, err
	}

	for _, b := range bookings {
		br := response.BookingResponse{
			ID:         b.ID,
			TotalPrice: b.TotalPrice,
			Status:     b.Status,
			CreatedAt:  b.CreatedAt,
			UpdatedAt:  b.UpdatedAt,
		}

		b.User.Password = ""

		if b.User.ID != uuid.Nil {
			br.User = b.User
		}

		if b.Match.ID != 0 {
			br.Match = b.Match
		}

		for _, bs := range b.Seats {
			br.Seats = append(br.Seats, response.BookingSeatResp{
				SeatID: bs.Seat.ID,
				SeatNo: bs.SeatNo,
				Price:  bs.Price,
			})
		}

		resp = append(resp, br)
	}

	return resp, http.StatusOK, nil
}

// update booking status
func (bs *bookingServiceImpl) UpdateBookingStatus(userID uuid.UUID, id uuid.UUID) (statusCode int, err error) {
	booking, err := bs.bookingRepository.GetByIDWithRelations(id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return http.StatusNotFound, errors.New("booking not found")
		}
		return http.StatusInternalServerError, err
	}

	if userID != booking.User.ID {
		// fmt.Println("userId : ", userID)
		// fmt.Println("booking user id : ", booking.User.ID)
		return http.StatusBadRequest, errors.New("booking cannot update by another user")
	}

	if booking.Status == "CANCELED" {
		return http.StatusBadRequest, errors.New("booking already cancelled")
	}

	if booking.Status == "PAID" {
		return http.StatusBadRequest, errors.New("booking already paid")
	}

	booking.Status = "PAID"

	if err := bs.bookingRepository.Update(booking); err != nil {
		return http.StatusInternalServerError, err
	}

	go bs.sendPaymentConfirmationEmail(booking)

	return http.StatusOK, nil
}

func (bs *bookingServiceImpl) sendPaymentConfirmationEmail(booking *Booking) {
	from := "stadiumhubtest@gmail.com"
	to := []string{booking.User.Email}
	appPassword := internal.ENV.EmailKey.EmailKey
	host := "smtp.gmail.com"
	addr := host + ":587"

	auth := smtp.PlainAuth("", from, strings.ReplaceAll(appPassword, " ", ""), host)

	// สร้างรายการที่นั่ง
	seatList := ""
	for _, bs := range booking.Seats {
		seatList += fmt.Sprintf("  - Seat No: %s (Price: $%.2f)\n", bs.SeatNo, bs.Price)
	}

	// สร้าง email message
	msg := ""
	msg += fmt.Sprintf("From: Stadium Hub <%s>\r\n", from)
	msg += fmt.Sprintf("To: %s\r\n", strings.Join(to, ","))
	msg += "Subject: Payment Confirmed - Booking Complete\r\n"
	msg += "MIME-Version: 1.0\r\n"
	msg += "Content-Type: text/plain; charset=UTF-8\r\n"
	msg += "Content-Transfer-Encoding: 8bit\r\n\r\n"

	msg += "Dear Customer,\n\n"
	msg += "✓ Payment Successful!\n\n"
	msg += "Your payment has been confirmed and your booking is now complete.\n\n"
	msg += "=== BOOKING DETAILS ===\n"
	msg += fmt.Sprintf("Booking ID: %s\n", booking.ID)
	msg += fmt.Sprintf("Match: %s vs %s\n", booking.Match.HomeTeam.Name, booking.Match.AwayTeam.Name)
	msg += fmt.Sprintf("Match Date: %s\n", booking.Match.UTCDate.Format("January 02, 2006 15:04 MST"))
	msg += fmt.Sprintf("Stadium: %s\n", booking.Match.HomeTeam.Venue) // ถ้ามีข้อมูล stadium
	msg += fmt.Sprintf("Total Paid: $%.2f\n\n", booking.TotalPrice)

	msg += "=== YOUR SEATS ===\n"
	msg += seatList
	msg += "\n"

	msg += "=== IMPORTANT INFORMATION ===\n"
	msg += "• Please arrive at least 30 minutes before the match starts\n"
	msg += "• Bring a valid ID for verification\n"
	msg += "• You can show this email or your booking ID at the entrance\n"
	msg += "• Your seats are reserved and guaranteed\n\n"

	msg += "=== MATCH DAY TIPS ===\n"
	msg += "• Gates open 1 hour before kickoff\n"
	msg += "• Food and beverages are available at the stadium\n"
	msg += "• Check weather conditions before you leave\n\n"

	msg += "We hope you enjoy the match!\n\n"
	msg += "If you have any questions, please don't hesitate to contact us.\n\n"
	msg += "Best regards,\n"
	msg += "Stadium Hub Team"

	// ส่งอีเมล
	if err := smtp.SendMail(addr, auth, from, to, []byte(msg)); err != nil {
		fmt.Printf("Failed to send payment confirmation email: %v\n", err)
	}
}

func (bs *bookingServiceImpl) GetRevenueByYear(year int) (resp []response.MonthRevenue, err error) {
	return bs.bookingRepository.GetRevenueByYear(year)
}

func (bs *bookingServiceImpl) CancelExpiredBookings(expireDuration time.Duration) (rows int, err error) {
	expireTime := time.Now().Add(-expireDuration)

	tx := bs.bookingRepository.Begin()
	if tx.Error != nil {
		return 0, tx.Error
	}

	result := tx.Model(&Booking{}).
		Where("status = ? AND created_at < ?", "PENDING", expireTime).
		Update("status", "CANCELED")

	if result.Error != nil {
		tx.Rollback()
		return 0, result.Error
	}

	if err := tx.Commit().Error; err != nil {
		return 0, err
	}

	return int(result.RowsAffected), nil
}
