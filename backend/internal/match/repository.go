package match

import (
	"fmt"
	"net/http"
	"time"

	"github.com/baimhons/stadiumhub/internal/team"
	"github.com/baimhons/stadiumhub/internal/utils"
	"gorm.io/gorm"
)

type MatchRepository interface {
	utils.BaseRepository[Match]
	GetTeamByID(teamID int) (*team.Team, error)
	UpdateOrCreateMatch(entity *Match) error
	GetAllMatches(query *utils.PaginationQuery) ([]Match, int, error)
	GetMatchesByTeamID(teamID uint, query *utils.PaginationQuery) ([]Match, int, error)
	GetMatchesByDateRange(startDate, endDate string, query *utils.PaginationQuery) ([]Match, int, error)
}

type matchRepositoryImpl struct {
	utils.BaseRepository[Match]
	db *gorm.DB
}

func NewMatchRepository(db *gorm.DB) MatchRepository {
	return &matchRepositoryImpl{
		BaseRepository: utils.NewBaseRepository[Match](db),
		db:             db,
	}
}

func (mr *matchRepositoryImpl) GetTeamByID(teamID int) (*team.Team, error) {
	var t team.Team
	if err := mr.db.Where("id = ?", teamID).First(&t).Error; err != nil {
		return nil, err
	}
	return &t, nil
}

func (mr *matchRepositoryImpl) UpdateOrCreateMatch(entity *Match) error {
	return mr.db.Save(entity).Error
}

func (mr *matchRepositoryImpl) GetAllMatches(query *utils.PaginationQuery) ([]Match, int, error) {
	var matches []Match
	tx := mr.db.Model(&Match{}).Preload("HomeTeam").Preload("AwayTeam")

	if query.Page != nil && query.PageSize != nil {
		offset := (*query.Page - 1) * (*query.PageSize)
		tx = tx.Offset(offset).Limit(*query.PageSize)
	}

	if query.Sort != nil && query.Order != nil {
		tx = tx.Order(*query.Sort + " " + *query.Order)
	}

	if err := tx.Find(&matches).Error; err != nil {
		return nil, http.StatusInternalServerError, err
	}

	if len(matches) == 0 {
		return nil, http.StatusNotFound, nil
	}

	return matches, http.StatusOK, nil
}

func (mr *matchRepositoryImpl) GetMatchesByTeamID(teamID uint, query *utils.PaginationQuery) ([]Match, int, error) {
	var matches []Match
	tx := mr.db.Model(&Match{}).
		Where("home_team_id = ? OR away_team_id = ?", teamID, teamID).
		Preload("HomeTeam").
		Preload("AwayTeam")

	// pagination
	if query.Page != nil && query.PageSize != nil {
		offset := (*query.Page - 1) * (*query.PageSize)
		tx = tx.Offset(offset).Limit(*query.PageSize)
	}

	// sort
	if query.Sort != nil && query.Order != nil {
		tx = tx.Order(*query.Sort + " " + *query.Order)
	}

	if err := tx.Find(&matches).Error; err != nil {
		return nil, http.StatusInternalServerError, err
	}

	if len(matches) == 0 {
		return nil, http.StatusNotFound, nil
	}

	return matches, http.StatusOK, nil
}

func (mr *matchRepositoryImpl) GetMatchesByDateRange(startDate, endDate string, query *utils.PaginationQuery) ([]Match, int, error) {
	var matches []Match

	// Parse startDate, endDate
	layout := "2006-01-02"
	start, err := time.Parse(layout, startDate)
	if err != nil {
		return nil, http.StatusBadRequest, fmt.Errorf("invalid startDate: %v", err)
	}
	end, err := time.Parse(layout, endDate)
	if err != nil {
		return nil, http.StatusBadRequest, fmt.Errorf("invalid endDate: %v", err)
	}

	// ดึงข้อมูลทั้งหมดก่อน filter วัน
	tx := mr.db.Model(&Match{}).Preload("HomeTeam").Preload("AwayTeam")

	// Sorting
	if query.Sort != nil && query.Order != nil {
		tx = tx.Order(*query.Sort + " " + *query.Order)
	}

	if err := tx.Find(&matches).Error; err != nil {
		return nil, http.StatusInternalServerError, err
	}

	// Filter matches ในช่วงวัน start-end (ตัดเวลาออก)
	var result []Match
	for _, m := range matches {
		y, mo, d := m.UTCDate.Date()
		startY, startMo, startD := start.Date()
		endY, endMo, endD := end.Date()

		matchDay := time.Date(y, mo, d, 0, 0, 0, 0, time.UTC)
		startDay := time.Date(startY, startMo, startD, 0, 0, 0, 0, time.UTC)
		endDay := time.Date(endY, endMo, endD, 0, 0, 0, 0, time.UTC)

		if !matchDay.Before(startDay) && !matchDay.After(endDay) {
			result = append(result, m)
		}
	}

	if len(result) == 0 {
		return nil, http.StatusNotFound, nil
	}

	// Pagination (ใช้ slice ของ result)
	page := 1
	if query.Page != nil && *query.Page > 0 {
		page = *query.Page
	}
	pageSize := 10
	if query.PageSize != nil && *query.PageSize > 0 {
		pageSize = *query.PageSize
	}

	startIndex := (page - 1) * pageSize
	endIndex := startIndex + pageSize
	if startIndex >= len(result) {
		return nil, http.StatusNotFound, nil
	}
	if endIndex > len(result) {
		endIndex = len(result)
	}

	return result[startIndex:endIndex], http.StatusOK, nil
}
