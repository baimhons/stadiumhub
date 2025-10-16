package response

import (
	"github.com/baimhons/stadiumhub/internal/team"
	"github.com/google/uuid"
)

type ZoneResponse struct {
	ID   uuid.UUID `json:"id"`
	Team team.Team `json:"team"`
	Name string    `json:"name"`
}
