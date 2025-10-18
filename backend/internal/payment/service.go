package payment

import (
	"net/http"
	"time"

	"github.com/baimhons/stadiumhub/internal"
	"github.com/baimhons/stadiumhub/internal/payment/api/response"
	"github.com/google/uuid"
	"github.com/stripe/stripe-go/v78"
	"github.com/stripe/stripe-go/v78/checkout/session"
)

type PaymentService interface {
	StripeSession(userID uuid.UUID, amount int64) (*response.PaymentResponse, error)
}

type paymentServiceImpl struct{}

func NewPaymentService() PaymentService {
	return &paymentServiceImpl{}
}

func (ps *paymentServiceImpl) StripeSession(userID uuid.UUID, amount int64) (*response.PaymentResponse, error) {
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
		StatusCode: http.StatusOK,
		SessionURL: stripeRecord,
	}

	return resp, nil
}
