package match

import (
	"time"

	"github.com/baimhons/stadiumhub/internal/team"
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
	ID         int       `gorm:"primaryKey"`
	HomeTeamID int       `gorm:"not null"`
	HomeTeam   team.Team `gorm:"foreignKey:HomeTeamID"`
	AwayTeamID int       `gorm:"not null"`
	AwayTeam   team.Team `gorm:"foreignKey:AwayTeamID"`
	UTCDate    time.Time `gorm:"not null"`
	Status     string    `gorm:"type:varchar(20);not null"`
	Venue      string    `gorm:"not null"`
	CreatedAt  time.Time `gorm:"not null default:current_timestamp"`
	UpdatedAt  time.Time `gorm:"not null default:current_timestamp"`
}
