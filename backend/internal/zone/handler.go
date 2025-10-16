package zone

import (
	"net/http"
	"strconv"

	"github.com/baimhons/stadiumhub/internal/utils"
	"github.com/gin-gonic/gin"
)

type ZoneHandler interface {
	GetZoneByMatchID(c *gin.Context)
}

type zoneHandlerImpl struct {
	zoneService ZoneService
}

func NewZoneHandler(zoneService ZoneService) ZoneHandler {
	return &zoneHandlerImpl{
		zoneService: zoneService,
	}
}
func (h *zoneHandlerImpl) GetZoneByMatchID(c *gin.Context) {
	idParam := c.Param("id")

	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusConflict, gin.H{"error": "invalid id"})
		return
	}

	resp, statusCode, err := h.zoneService.GetZoneByMatchID(id)
	if err != nil {
		c.JSON(statusCode, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, utils.SuccessResponse{
		Message: "get zone successfully",
		Data:    resp,
	})
}
