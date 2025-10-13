package seat

import (
	"net/http"
	"strconv"

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
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid match_id"})
		return
	}

	teamID, err := strconv.Atoi(teamIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid team_id"})
		return
	}

	var zoneID *uuid.UUID
	if zoneIDStr != "" {
		z, err := uuid.Parse(zoneIDStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid zone_id"})
			return
		}
		zoneID = &z
	}

	seats, err := h.seatService.GetAvailableSeats(uint(matchID), teamID, zoneID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, seats)
}
