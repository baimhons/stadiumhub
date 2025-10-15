package booking

import (
	"net/http"

	"github.com/baimhons/stadiumhub/internal/booking/api/request"
	"github.com/baimhons/stadiumhub/internal/utils"
	"github.com/gin-gonic/gin"
)

type BookingValidate interface {
	ValidateSeatQuantity(c *gin.Context)
}

type bookingValidate struct {
}

func NewBookingValidate() *bookingValidate {
	return &bookingValidate{}
}

func (b *bookingValidate) ValidateSeatQuantity(c *gin.Context) {
	var req request.CreateBookingRequest

	if err := utils.ValidateCommonRequestBody(c, &req); err != nil {
		c.JSON(http.StatusBadRequest, utils.ErrorResponse{
			Message: err.Error(),
		})
		c.Abort()
		return
	}

	if len(req.SeatIDs) > 6 {
		c.JSON(http.StatusBadRequest, utils.ErrorResponse{
			Message: "can only booking 6 seat",
		})
		c.Abort()
		return
	}

	c.Set("req", req)
	c.Next()
}
