package booking

import (
	"net/http"
	"strconv"

	"github.com/baimhons/stadiumhub/internal/booking/api/request"
	"github.com/baimhons/stadiumhub/internal/models"
	"github.com/baimhons/stadiumhub/internal/utils"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type BookingHandler interface {
	GetRevenueByYear(c *gin.Context)
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
		c.JSON(http.StatusBadRequest, utils.ErrorResponse{
			Message: "request not found",
			Error:   nil,
		})
		return
	}

	req, ok := reqRaw.(request.CreateBookingRequest)
	if !ok {
		c.JSON(http.StatusBadRequest, utils.ErrorResponse{
			Message: "invalid request data",
			Error:   nil,
		})
		return
	}

	userCtxRaw, exists := c.Get("userContext")
	if !exists {
		c.JSON(http.StatusUnauthorized, utils.ErrorResponse{
			Message: "unauthorized",
			Error:   nil,
		})
		return
	}
	userCtx, ok := userCtxRaw.(models.UserContext)
	if !ok {
		c.JSON(http.StatusBadRequest, utils.ErrorResponse{
			Message: "invalid user context",
			Error:   nil,
		})
		return
	}

	resp, statusCode, err := h.bookingService.CreateBooking(userCtx, req)
	if err != nil {
		c.JSON(statusCode, utils.ErrorResponse{
			Message: err.Error(),
			Error:   err,
		})
		return
	}

	c.JSON(statusCode, utils.SuccessResponse{
		Message: "create booking success",
		Data:    resp,
	})
}

func (h *bookingHandlerImpl) GetBookingByID(c *gin.Context) {
	idParam := c.Param("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, utils.ErrorResponse{
			Message: "invalid booking id",
			Error:   err,
		})
		return
	}

	userCtxRaw, exists := c.Get("userContext")
	if !exists {
		c.JSON(http.StatusUnauthorized, utils.ErrorResponse{
			Message: "unauthorized",
			Error:   nil,
		})
		return
	}
	userCtx, ok := userCtxRaw.(models.UserContext)
	if !ok {
		c.JSON(http.StatusBadRequest, utils.ErrorResponse{
			Message: "invalid user context",
			Error:   nil,
		})
		return
	}

	resp, statusCode, err := h.bookingService.GetBookingByID(id, userCtx)
	if err != nil {
		c.JSON(statusCode, utils.ErrorResponse{
			Message: err.Error(),
			Error:   err,
		})
		return
	}

	c.JSON(http.StatusOK, utils.SuccessResponse{
		Message: "get booking success",
		Data:    resp,
	})
}

func (h *bookingHandlerImpl) GetAllBookingsByUser(c *gin.Context) {
	userCtxRaw, exists := c.Get("userContext")
	if !exists {
		c.JSON(http.StatusUnauthorized, utils.ErrorResponse{
			Message: "unauthorized",
			Error:   nil,
		})
		return
	}
	userCtx, ok := userCtxRaw.(models.UserContext)
	if !ok {
		c.JSON(http.StatusBadRequest, utils.ErrorResponse{
			Message: "invalid user context",
			Error:   nil,
		})
		return
	}

	var query utils.PaginationQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		c.JSON(http.StatusBadRequest, utils.ErrorResponse{
			Message: err.Error(),
			Error:   err,
		})
		return
	}

	resp, statusCode, err := h.bookingService.GetAllBookingsByUser(userCtx.ID, &query)
	if err != nil {
		c.JSON(statusCode, utils.ErrorResponse{
			Message: err.Error(),
			Error:   err,
		})
		return
	}

	c.JSON(statusCode, utils.SuccessResponse{
		Message: "success",
		Data:    resp,
	})
}

func (h *bookingHandlerImpl) CancelBooking(c *gin.Context) {
	userCtxRaw, exists := c.Get("userContext")
	if !exists {
		c.JSON(http.StatusUnauthorized, utils.ErrorResponse{
			Message: "unauthorized",
			Error:   nil,
		})
		return
	}
	userCtx, ok := userCtxRaw.(models.UserContext)
	if !ok {
		c.JSON(http.StatusBadRequest, utils.ErrorResponse{
			Message: "invalid user context",
			Error:   nil,
		})
		return
	}
	idParam := c.Param("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, utils.ErrorResponse{
			Message: "invalid booking id",
			Error:   err,
		})
		return
	}

	statusCode, err := h.bookingService.CancelBooking(userCtx.ID, id)
	if err != nil {
		c.JSON(statusCode, utils.ErrorResponse{
			Message: err.Error(),
			Error:   err,
		})
		return
	}

	c.JSON(http.StatusOK, utils.SuccessResponse{
		Message: "booking cancelled successfully",
		Data:    nil,
	})
}

func (h *bookingHandlerImpl) GetAllBookings(c *gin.Context) {
	var query utils.PaginationQuery

	if err := c.ShouldBindQuery(&query); err != nil {
		c.JSON(http.StatusBadRequest, utils.ErrorResponse{
			Message: err.Error(),
			Error:   err,
		})
		return
	}

	resp, statusCode, err := h.bookingService.GetAllBookings(&query)
	if err != nil {
		c.JSON(statusCode, utils.ErrorResponse{
			Message: err.Error(),
			Error:   err,
		})
		return
	}

	c.JSON(http.StatusOK, utils.SuccessResponse{
		Message: "success",
		Data:    resp,
	})
}

func (h *bookingHandlerImpl) UpdateBookingStatus(c *gin.Context) {
	userCtxRaw, exists := c.Get("userContext")
	if !exists {
		c.JSON(http.StatusUnauthorized, utils.ErrorResponse{
			Message: "unauthorized",
			Error:   nil,
		})
		return
	}
	userCtx, ok := userCtxRaw.(models.UserContext)
	if !ok {
		c.JSON(http.StatusBadRequest, utils.ErrorResponse{
			Message: "invalid user context",
			Error:   nil,
		})
		return
	}
	idParam := c.Param("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, utils.ErrorResponse{
			Message: "invalid booking id",
			Error:   err,
		})
		return
	}

	statusCode, err := h.bookingService.UpdateBookingStatus(userCtx.ID, id)
	if err != nil {
		c.JSON(statusCode, utils.ErrorResponse{
			Message: err.Error(),
			Error:   err,
		})
		return
	}

	c.JSON(http.StatusOK, utils.SuccessResponse{
		Message: "booking update successfully",
		Data:    nil,
	})
}

func (h *bookingHandlerImpl) GetRevenueByYear(c *gin.Context) {
	yearStr := c.Query("year")
	if yearStr == "" {
		c.JSON(http.StatusBadRequest, utils.ErrorResponse{
			Message: "Year is required",
		})
		return
	}

	year, err := strconv.Atoi(yearStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, utils.ErrorResponse{
			Message: "Invalid year",
		})
		return
	}

	revenueMap, err := h.bookingService.GetRevenueByYear(year)
	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.ErrorResponse{
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, utils.SuccessResponse{
		Message: "Revenue by month successfully",
		Data:    revenueMap,
	})
}
