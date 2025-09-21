package matchteam

import (
	"github.com/baimhons/stadiumhub.git/internal/match"
	"github.com/baimhons/stadiumhub.git/internal/team"
	"github.com/baimhons/stadiumhub.git/internal/utils"
	"github.com/google/uuid"
)

type MatchTeam struct {
	utils.BaseEntity
	MatchID uuid.UUID
	Match   match.Match `gorm:"foreignkey:MatchID"`
	TeamID  uuid.UUID
	Team    team.Team `gorm:"foreignkey:TeamID"`
}
