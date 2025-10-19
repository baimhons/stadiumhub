package payment

import (
	"errors"
	"net/http"
	"time"

	"github.com/baimhons/stadiumhub/internal"
	"github.com/baimhons/stadiumhub/internal/booking"
	"github.com/baimhons/stadiumhub/internal/payment/api/response"
	"github.com/google/uuid"
	"github.com/stripe/stripe-go/v78"
	"github.com/stripe/stripe-go/v78/checkout/session"
	"gorm.io/gorm"
)

type PaymentService interface {
	StripeSession(userID uuid.UUID, bookingID uuid.UUID, amount int64) (*response.PaymentResponse, error)
}

type paymentServiceImpl struct {
	bookingRepository booking.BookingRepository
}

func NewPaymentService(bookingRepository booking.BookingRepository) PaymentService {
	return &paymentServiceImpl{
		bookingRepository: bookingRepository,
	}
}

func (ps *paymentServiceImpl) StripeSession(userID uuid.UUID, bookingID uuid.UUID, amount int64) (*response.PaymentResponse, error) {

	if err := ps.bookingRepository.GetByID(&booking.Booking{}, bookingID); err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.New("booking not found")
		}
		return nil, err
	}

	// if booking.Booking.Status == "PENDING" {

	// }

	stripe.Key = internal.ENV.Stripe.StripeKey

	amount = amount * 100

	params := &stripe.CheckoutSessionParams{
		Mode: stripe.String(string(stripe.CheckoutSessionModePayment)),
		LineItems: []*stripe.CheckoutSessionLineItemParams{{
			PriceData: &stripe.CheckoutSessionLineItemPriceDataParams{
				Currency: stripe.String("thb"),
				ProductData: &stripe.CheckoutSessionLineItemPriceDataProductDataParams{
					Name: stripe.String("Football Ticket"),
				},
				UnitAmount: stripe.Int64(amount),
			},
			Quantity: stripe.Int64(1),
		}},
		SuccessURL: stripe.String("http://127.0.0.1:5500/frontend/pages/payment/success.html"),
		CancelURL:  stripe.String("http://127.0.0.1:5500/frontend/pages/payment/cancel.html"),
		ExpiresAt:  stripe.Int64(time.Now().Add(30 * time.Minute).Unix()),
	}

	s, err := session.New(params)
	if err != nil {
		return nil, err
	}

	stripeRecord := response.StripeModel{
		UserID:     userID,
		SessionURL: s.URL,
		SessionID:  s.ID,
		Amount:     float32(amount) / 100,
		CreatedAt:  time.Now(),
	}

	resp := &response.PaymentResponse{
		Status:     1,
		Message:    "Checkout Session Successfully!",
		BookingID:  bookingID,
		StatusCode: http.StatusOK,
		SessionURL: stripeRecord,
	}

	return resp, nil
}
