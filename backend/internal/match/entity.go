package match

import (
	"time"

	"github.com/baimhons/stadiumhub/internal/team"
	"github.com/baimhons/stadiumhub/internal/utils"
	"github.com/google/uuid"
)

// type MatchStatus string

// const (
// 	MatchScheduled MatchStatus = "SCHEDULED"
// 	MatchLive      MatchStatus = "LIVE"
// 	MatchFinished  MatchStatus = "FINISHED"
// 	MatchPostponed MatchStatus = "POSTPONED"
// 	MatchInPlay    MatchStatus = "IN_PLAY"
// 	MatchPaused    MatchStatus = "PAUSED"
// 	MatchSuspended MatchStatus = "SUSPENDED"
// 	MatchCancelled MatchStatus = "CANCELLED"
// )

type Match struct {
	utils.BaseEntity
	HomeTeamID uuid.UUID `gorm:"not null"`
	HomeTeam   team.Team `gorm:"foreignKey:HomeTeamID"`
	AwayTeamID uuid.UUID `gorm:"not null"`
	AwayTeam   team.Team `gorm:"foreignKey:AwayTeamID"`
	UTCDate    time.Time `gorm:"not null"`
	Status     string    `gorm:"type:varchar(20);not null"`
}
