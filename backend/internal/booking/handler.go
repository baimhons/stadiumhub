package booking

import (
	"net/http"

	"github.com/baimhons/stadiumhub/internal/booking/api/request"
	"github.com/baimhons/stadiumhub/internal/models"
	"github.com/gin-gonic/gin"
)

type BookingHandler interface {
	CreateBooking(c *gin.Context)
}

type bookingHandlerImpl struct {
	bookingService BookingService
}

func NewBookingHandler(bookingService BookingService) BookingHandler {
	return &bookingHandlerImpl{
		bookingService: bookingService,
	}
}

func (h *bookingHandlerImpl) CreateBooking(c *gin.Context) {
	reqRaw, exists := c.Get("req")
	if !exists {
		c.JSON(http.StatusBadRequest, gin.H{"message": "request not found"})
		return
	}

	req, ok := reqRaw.(request.CreateBookingRequest)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"message": "invalid request data"})
		return
	}

	userCtxRaw, exists := c.Get("userContext")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "unauthorized"})
		return
	}

	userCtx, ok := userCtxRaw.(models.UserContext)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"message": "invalid user context"})
		return
	}

	resp, statusCode, err := h.bookingService.CreateBooking(userCtx, req)
	if err != nil {
		c.JSON(statusCode, gin.H{"message": err.Error()})
		return
	}

	c.JSON(statusCode, resp)
}
