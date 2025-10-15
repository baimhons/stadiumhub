package initial

import (
	"github.com/baimhons/stadiumhub/internal/booking"
	"github.com/baimhons/stadiumhub/internal/user"
)

type Validate struct {
	UserValidate    user.UserValidate
	BookingValidate booking.BookingValidate
}

func NewValidate() *Validate {
	return &Validate{
		UserValidate:    user.NewUserValidate(),
		BookingValidate: booking.NewBookingValidate(),
	}
}
