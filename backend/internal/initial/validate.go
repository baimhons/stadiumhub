package initial

import (
	"github.com/baimhons/stadiumhub/internal/booking"
	"github.com/baimhons/stadiumhub/internal/payment"
	"github.com/baimhons/stadiumhub/internal/user"
)

type Validate struct {
	UserValidate    user.UserValidate
	BookingValidate booking.BookingValidate
	PaymentValidate payment.PaymentValidate
}

func NewValidate() *Validate {
	return &Validate{
		UserValidate:    user.NewUserValidate(),
		BookingValidate: booking.NewBookingValidate(),
		PaymentValidate: payment.NewPaymentValidate(),
	}
}
