package team

import (
	"net/http"

	"github.com/baimhons/stadiumhub/internal/utils"
	"github.com/gin-gonic/gin"
)

type TeamHandler interface {
	GetAllTeam(c *gin.Context)
}

type teamHandlerImpl struct {
	teamService TeamService
}

func NewTeamHandler(teamService TeamService) TeamHandler {
	return &teamHandlerImpl{
		teamService: teamService,
	}
}

func (h *teamHandlerImpl) GetAllTeam(c *gin.Context) {
	query := utils.PaginationQuery{}

	if err := c.ShouldBindQuery(&query); err != nil {
		c.JSON(http.StatusBadRequest, utils.ErrorResponse{
			Message: "invalid query parameters",
			Error:   err,
		})
		return
	}
	defaultSize := 20
	query.PageSize = &defaultSize

	defaultPage := 0
	query.Page = &defaultPage

	teams, err := h.teamService.GetAllTeam(&query)
	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.ErrorResponse{
			Message: err.Error(),
			Error:   err,
		})
		return
	}

	c.JSON(http.StatusOK, utils.SuccessResponse{
		Message: "get teams successfully",
		Data:    teams,
	})
}
