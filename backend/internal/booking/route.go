package booking

import (
	"github.com/baimhons/stadiumhub/internal/middlewares"
	"github.com/baimhons/stadiumhub/internal/user"
	"github.com/gin-gonic/gin"
)

type BookingRoutes struct {
	group           *gin.RouterGroup
	bookingHandler  BookingHandler
	bookingValidate BookingValidate
	userValidate    user.UserValidate
	authMiddleware  middlewares.AuthMiddleware
}

func NewBookingRoutes(
	group *gin.RouterGroup,
	bookingHandler BookingHandler,
	bookingValidate BookingValidate,
	userValidate user.UserValidate,
	authMiddleware middlewares.AuthMiddleware,
) *BookingRoutes {

	bookingGroup := group.Group("/booking")
	r := &BookingRoutes{
		group:           bookingGroup,
		bookingHandler:  bookingHandler,
		bookingValidate: bookingValidate,
		userValidate:    userValidate,
		authMiddleware:  authMiddleware,
	}

	return r
}

func (r *BookingRoutes) RegisterRoutes() {

	r.group.GET("/revenue", r.authMiddleware.RequireAuth(), r.userValidate.ValidateRoleAdmin, r.bookingHandler.GetRevenueByYear)
	r.group.GET("/:id", r.authMiddleware.RequireAuth(), r.bookingHandler.GetBookingByID)
	r.group.GET("/history", r.authMiddleware.RequireAuth(), r.bookingHandler.GetAllBookingsByUser)
	r.group.GET("/all", r.authMiddleware.RequireAuth(), r.userValidate.ValidateRoleAdmin, r.bookingHandler.GetAllBookings)
	r.group.POST("/create", r.authMiddleware.RequireAuth(), r.bookingValidate.ValidateSeatQuantity, r.bookingHandler.CreateBooking)
	r.group.POST("/update-status/:id", r.authMiddleware.RequireAuth(), r.bookingHandler.UpdateBookingStatus)
	r.group.POST("/cancel/:id", r.authMiddleware.RequireAuth(), r.bookingHandler.CancelBooking)
}
