package seed

import (
	"fmt"

	"github.com/baimhons/stadiumhub/internal/team"
	"github.com/baimhons/stadiumhub/internal/zone"
	"gorm.io/gorm"
)

func SeedZones(db *gorm.DB) error {
	zoneNames := []string{"north", "south", "east", "west"}

	var teams []team.Team
	if err := db.Find(&teams).Error; err != nil {
		return fmt.Errorf("failed to get teams: %w", err)
	}

	for _, team := range teams {
		for _, zoneName := range zoneNames {
			z := zone.Zone{
				TeamID: team.ID,
				Name:   zoneName,
			}

			if err := db.FirstOrCreate(&z, zone.Zone{TeamID: team.ID, Name: zoneName}).Error; err != nil {
				return fmt.Errorf("failed to create zone for team %d (%s): %w", team.ID, team.Name, err)
			}
		}
	}

	return nil
}
