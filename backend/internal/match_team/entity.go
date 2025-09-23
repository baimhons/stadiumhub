package matchteam

import (
	"github.com/baimhons/stadiumhub/internal/match"
	"github.com/baimhons/stadiumhub/internal/team"
	"github.com/baimhons/stadiumhub/internal/utils"
	"github.com/google/uuid"
)

type MatchTeam struct {
	utils.BaseEntity
	MatchID uuid.UUID
	Match   match.Match `gorm:"foreignkey:MatchID"`
	TeamID  uuid.UUID
	Team    team.Team `gorm:"foreignkey:TeamID"`
}
