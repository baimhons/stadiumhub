package match

import (
	"net/http"
	"strconv"
	"time"

	"github.com/baimhons/stadiumhub/internal/utils"
	"github.com/gin-gonic/gin"
)

type MatchHandler interface {
	UpdateMatches(c *gin.Context)
	GetAllMatches(c *gin.Context)
	GetMatchesByTeamID(c *gin.Context)
	GetMatchesByDateRange(c *gin.Context)
}

type matchHandlerImpl struct {
	matchService MatchService
}

func NewMatchHandler(matchService MatchService) MatchHandler {
	return &matchHandlerImpl{matchService: matchService}
}

func (h *matchHandlerImpl) UpdateMatches(c *gin.Context) {
	monthStr := c.Query("month")
	yearStr := c.Query("year")

	now := time.Now()
	month := int(now.Month())
	year := now.Year()

	if monthStr != "" {
		if m, err := strconv.Atoi(monthStr); err == nil {
			month = m
		}
	}
	if yearStr != "" {
		if y, err := strconv.Atoi(yearStr); err == nil {
			year = y
		}
	}

	msg, status, err := h.matchService.UpdateMatches(month, year)
	if err != nil {
		c.JSON(status, gin.H{"error": err.Error()})
		return
	}
	c.JSON(status, gin.H{"message": msg})
}

func (h *matchHandlerImpl) GetAllMatches(c *gin.Context) {
	var query utils.PaginationQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	data, status, err := h.matchService.GetAllMatches(&query)
	if err != nil {
		c.JSON(status, gin.H{"error": err.Error()})
		return
	}
	c.JSON(status, data)
}

func (h *matchHandlerImpl) GetMatchesByTeamID(c *gin.Context) {
	var query utils.PaginationQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	teamID, _ := strconv.ParseUint(c.Param("teamID"), 10, 64)
	data, status, err := h.matchService.GetMatchesByTeamID(uint(teamID), &query)
	if err != nil {
		c.JSON(status, gin.H{"error": err.Error()})
		return
	}
	c.JSON(status, data)
}

func (h *matchHandlerImpl) GetMatchesByDateRange(c *gin.Context) {
	var query utils.PaginationQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	startDate := c.Query("startDate")
	endDate := c.Query("endDate")

	data, status, err := h.matchService.GetMatchesByDateRange(startDate, endDate, &query)
	if err != nil {
		c.JSON(status, gin.H{"error": err.Error()})
		return
	}
	c.JSON(status, data)
}
