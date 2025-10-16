package seat

import (
	"net/http"
	"strconv"

	"github.com/baimhons/stadiumhub/internal/utils"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type SeatHandler interface {
	GetAvailableSeats(c *gin.Context)
}

type seatHandlerImpl struct {
	seatService SeatService
}

func NewSeatHandler(seatService SeatService) SeatHandler {
	return &seatHandlerImpl{seatService: seatService}
}

func (h *seatHandlerImpl) GetAvailableSeats(c *gin.Context) {
	matchIDStr := c.Query("match_id")
	teamIDStr := c.Query("team_id")
	zoneIDStr := c.Query("zone_id")

	matchID, err := strconv.ParseUint(matchIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, utils.ErrorResponse{
			Message: "invalid match_id",
			Error:   nil,
		})
		return
	}

	teamID, err := strconv.Atoi(teamIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, utils.ErrorResponse{
			Message: "invalid team_id",
			Error:   err,
		})
		return
	}

	var zoneID *uuid.UUID
	if zoneIDStr != "" {
		z, err := uuid.Parse(zoneIDStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, utils.ErrorResponse{
				Message: "invalid zone_id",
				Error:   err,
			})
			return
		}
		zoneID = &z
	}

	seats, err := h.seatService.GetAvailableSeats(uint(matchID), teamID, zoneID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.ErrorResponse{
			Message: err.Error(),
			Error:   err,
		})
		return
	}

	c.JSON(http.StatusOK, seats)
}
