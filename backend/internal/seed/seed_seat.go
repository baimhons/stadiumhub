package seed

import (
	"fmt"

	"github.com/baimhons/stadiumhub/internal/seat"
	"github.com/baimhons/stadiumhub/internal/team"
	"github.com/baimhons/stadiumhub/internal/zone"
	"gorm.io/gorm"
)

func SeedSeats(db *gorm.DB) error {
	var teams []team.Team
	if err := db.Find(&teams).Error; err != nil {
		return fmt.Errorf("failed to get teams: %w", err)
	}

	for _, team := range teams {
		var zones []zone.Zone
		if err := db.Where("team_id = ?", team.ID).Find(&zones).Error; err != nil {
			return fmt.Errorf("failed to get zones for team %d: %w", team.ID, err)
		}

		if len(zones) == 0 {
			continue
		}

		seatPerZone := (team.ViewerCapacity / 100) / len(zones)
		remainder := (team.ViewerCapacity / 100) % len(zones)

		for i, zone := range zones {
			count := seatPerZone
			if i < remainder {
				count++
			}

			for j := 1; j <= count; j++ {
				s := seat.Seat{
					SeatNo: fmt.Sprintf("%s-%d", zone.Name[:1], j), // N-1, S-1, E-1, W-1
					ZoneID: zone.ID,
				}

				if err := db.FirstOrCreate(&s, seat.Seat{
					SeatNo: s.SeatNo,
					ZoneID: zone.ID,
				}).Error; err != nil {
					return fmt.Errorf("failed to create seat %s: %w", s.SeatNo, err)
				}
			}

			fmt.Printf("Created %d seats for Team %s Zone %s\n", count, team.Name, zone.Name)
		}
	}

	return nil
}
