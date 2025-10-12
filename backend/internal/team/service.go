package team

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/baimhons/stadiumhub/internal/utils"
)

type TeamService interface {
	InsertTeamCapacityAndPrice() error
}
type teamServiceImpl struct {
	teamRepository TeamRepository
}

func NewTeamService(teamRepo TeamRepository) TeamService {
	return &teamServiceImpl{teamRepository: teamRepo}
}

type StadiumInfo struct {
	Team     string  `json:"team"`
	Stadium  string  `json:"stadium"`
	Capacity int     `json:"capacity"`
	Price    float32 `json:"price"`
}

func (ts *teamServiceImpl) InsertTeamCapacityAndPrice() error {
	filePath := "C:/Users/gigam/repository/project_web_dev/backend/internal/data/premier_league_stadiums.json"
	data, err := os.ReadFile(filePath)

	if err != nil {
		return fmt.Errorf("error reading file: %w", err)
	}

	var stadiums []StadiumInfo
	if err := json.Unmarshal(data, &stadiums); err != nil {
		return fmt.Errorf("error unmarshalling json: %w", err)
	}

	teams := []Team{}

	page := 0
	pageSize := 30
	sortField := "name"
	order := "asc"

	query := &utils.PaginationQuery{
		Page:     &page,
		PageSize: &pageSize,
		Sort:     &sortField,
		Order:    &order,
	}

	if err := ts.teamRepository.GetAll(&teams, query); err != nil {
		return fmt.Errorf("error fetching teams: %w", err)
	}

	type StadiumData struct {
		Capacity int
		Price    float32
	}
	stadiumMap := make(map[string]StadiumData)
	for _, s := range stadiums {
		stadiumMap[s.Team] = StadiumData{
			Capacity: s.Capacity,
			Price:    s.Price,
		}
	}

	for i := range teams {
		if stadiumData, ok := stadiumMap[teams[i].ShortName]; ok {
			teams[i].ViewerCapacity = stadiumData.Capacity
			teams[i].Price = stadiumData.Price
			if err := ts.teamRepository.Update(&teams[i]); err != nil {
				return fmt.Errorf("error updating team %s: %w", teams[i].Name, err)
			}
		}
	}

	return nil
}
