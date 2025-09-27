package match

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/baimhons/stadiumhub/internal"
	"github.com/baimhons/stadiumhub/internal/match/api/response"
	"github.com/baimhons/stadiumhub/internal/team"
	"github.com/baimhons/stadiumhub/internal/utils"
	"gorm.io/gorm"
)

func fetchMatches(apiKey, dateFrom, dateTo string) ([]response.ApiMatch, error) {
	url := fmt.Sprintf(
		"https://api.football-data.org/v4/matches?competitions=PL&dateFrom=%s&dateTo=%s",
		dateFrom, dateTo,
	)

	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("X-Auth-Token", apiKey)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var result response.MatchResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, err
	}

	return result.Matches, nil
}

func SeedMatches(db *gorm.DB) {
	apiKey := internal.ENV.APIFootball.APIKey

	now := time.Now()
	firstDay := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, time.UTC)
	lastDay := firstDay.AddDate(0, 1, -1)

	matchRepo := utils.NewBaseRepository[Match](db)

	for start := firstDay; start.Before(lastDay); start = start.AddDate(0, 0, 7) {
		end := start.AddDate(0, 0, 6)
		if end.After(lastDay) {
			end = lastDay
		}

		apiMatches, err := fetchMatches(apiKey, start.Format("2006-01-02"), end.Format("2006-01-02"))
		if err != nil {
			panic(err)
		}

		for _, m := range apiMatches {
			utcDate, _ := time.Parse(time.RFC3339, m.UtcDate)

			utcDate = utcDate.In(time.UTC)
			t := team.Team{}
			if err := db.Where("ID = ?", m.HomeTeam.ID).First(&t).Error; err != nil {
				fmt.Printf("❌ team not found: %d\n", m.HomeTeam.ID)
				continue
			}
			venue := ifEmpty(t.Venue, "Unknown Stadium")

			entity := Match{
				ID:         m.ID,
				UTCDate:    utcDate,
				Status:     m.Status,
				HomeTeamID: m.HomeTeam.ID,
				AwayTeamID: m.AwayTeam.ID,
				Venue:      venue,
			}

			if err := matchRepo.Create(&entity); err != nil {
				fmt.Printf("❌ save failed match %d: %v\n", entity.ID, err)
			} else {
				fmt.Printf("✅ saved match: %d vs %d (%s)\n",
					entity.HomeTeam.ID, entity.AwayTeam.ID, entity.UTCDate.Format("2006-01-02"))
			}
		}
	}
}

func ifEmpty(val, fallback string) string {
	if val == "" || val == "null" {
		return fallback
	}
	return val
}
