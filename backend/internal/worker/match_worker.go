package worker

import (
	"fmt"
	"time"

	"github.com/baimhons/stadiumhub/internal/match"
)

type MatchWorker struct {
	matchService match.MatchService
}

func NewMatchWorker(matchService match.MatchService) *MatchWorker {
	return &MatchWorker{matchService: matchService}
}

func (w *MatchWorker) Start() {
	ticker := time.NewTicker(time.Hour)
	defer ticker.Stop()

	for range ticker.C {
		now := time.Now()
		month := int(now.Month())
		year := now.Year()

		fmt.Println("[Worker] Auto updating matches...")
		msg, _, err := w.matchService.UpdateMatches(month, year)
		if err != nil {
			fmt.Println("[Worker] Failed:", err)
		} else {
			fmt.Println("[Worker]", msg)
		}
	}
}
