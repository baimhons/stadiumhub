package booking

import (
	"net/http"

	"github.com/baimhons/stadiumhub/internal/booking/api/request"
	"github.com/baimhons/stadiumhub/internal/models"
	"github.com/baimhons/stadiumhub/internal/utils"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type BookingHandler interface {
	CreateBooking(c *gin.Context)
	GetBookingByID(c *gin.Context)
	GetAllBookingsByUser(c *gin.Context)
	CancelBooking(c *gin.Context)
	GetAllBookings(c *gin.Context)
	UpdateBookingStatus(c *gin.Context)
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

func (h *bookingHandlerImpl) GetBookingByID(c *gin.Context) {
	idParam := c.Param("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "invalid booking id"})
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

	resp, statusCode, err := h.bookingService.GetBookingByID(id, userCtx)
	if err != nil {
		c.JSON(statusCode, gin.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, resp)
}

func (h *bookingHandlerImpl) GetAllBookingsByUser(c *gin.Context) {
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

	var query utils.PaginationQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	resp, statusCode, err := h.bookingService.GetAllBookingsByUser(userCtx.ID, &query)
	if err != nil {
		c.JSON(statusCode, gin.H{"error": err.Error()})
		return
	}

	c.JSON(statusCode, gin.H{
		"message": "success",
		"data":    resp,
	})
}

func (h *bookingHandlerImpl) CancelBooking(c *gin.Context) {
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
	idParam := c.Param("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid booking id"})
		return
	}

	statusCode, err := h.bookingService.CancelBooking(userCtx.ID, id)
	if err != nil {
		c.JSON(statusCode, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "booking cancelled successfully",
	})
}

func (h *bookingHandlerImpl) GetAllBookings(c *gin.Context) {
	var query utils.PaginationQuery

	if err := c.ShouldBindQuery(&query); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	resp, statusCode, err := h.bookingService.GetAllBookings(&query)
	if err != nil {
		c.JSON(statusCode, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "success",
		"data":    resp,
	})
}

func (h *bookingHandlerImpl) UpdateBookingStatus(c *gin.Context) {
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
	idParam := c.Param("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid booking id"})
		return
	}

	statusCode, err := h.bookingService.UpdateBookingStatus(userCtx.ID, id)
	if err != nil {
		c.JSON(statusCode, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "booking update successfully",
	})
}
