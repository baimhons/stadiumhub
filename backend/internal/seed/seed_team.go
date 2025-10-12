package seed

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/baimhons/stadiumhub/internal"
	"github.com/baimhons/stadiumhub/internal/team"
	"gorm.io/gorm"
)

type ApiTeam struct {
	ID        int    `json:"id"`
	Name      string `json:"name"`
	ShortName string `json:"shortName"`
	TLA       string `json:"tla"`
	Address   string `json:"address"`
	Venue     string `json:"venue"`
}

type TeamResponse struct {
	Teams []ApiTeam `json:"teams"`
}

func SeedTeam(db *gorm.DB) {
	apiKey := internal.ENV.APIFootball.APIKey
	url := "https://api.football-data.org/v4/competitions/PL/teams"

	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("X-Auth-Token", apiKey)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	var result TeamResponse
	if err := json.Unmarshal(body, &result); err != nil {
		panic(err)
	}

	teamRepo := team.NewTeamRepository(db)
	if teamRepo == nil {
		panic("team repository is nil")
	}

	for _, apiTeam := range result.Teams {
		entity := team.Team{
			ID:             apiTeam.ID,
			Name:           apiTeam.Name,
			ShortName:      apiTeam.ShortName,
			TLA:            apiTeam.TLA,
			Address:        apiTeam.Address,
			Venue:          apiTeam.Venue,
			ViewerCapacity: 0,
			Price:          0.0,
		}

		if err := teamRepo.Create(&entity); err != nil {
			fmt.Printf("error saving team %s: %v\n", entity.Name, err)
		}
	}
	service := team.NewTeamService(teamRepo)
	if err := service.InsertTeamCapacityAndPrice(); err != nil {
		fmt.Printf("error inserting team capacity: %v\n", err)
	}
}
