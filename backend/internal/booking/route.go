package booking

import (
	"github.com/baimhons/stadiumhub/internal/middlewares"
	"github.com/gin-gonic/gin"
)

type BookingRoutes struct {
	group           *gin.RouterGroup
	bookingHandler  BookingHandler
	bookingValidate BookingValidate
	authMiddleware  middlewares.AuthMiddleware
}

func NewBookingRoutes(
	group *gin.RouterGroup,
	bookingHandler BookingHandler,
	bookingValidate BookingValidate,
	authMiddleware middlewares.AuthMiddleware,
) *BookingRoutes {

	bookingGroup := group.Group("/booking")
	r := &BookingRoutes{
		group:           bookingGroup,
		bookingHandler:  bookingHandler,
		bookingValidate: bookingValidate,
		authMiddleware:  authMiddleware,
	}

	return r
}

func (r *BookingRoutes) RegisterRoutes() {

	r.group.POST("/create", r.authMiddleware.RequireAuth(), r.bookingValidate.ValidateSeatQuantity, r.bookingHandler.CreateBooking)
}
