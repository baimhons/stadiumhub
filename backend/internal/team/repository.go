package team

import (
	"github.com/baimhons/stadiumhub/internal/utils"
	"gorm.io/gorm"
)

type TeamRepository interface {
	utils.BaseRepository[Team]
}

type teamRepositoryImpl struct {
	utils.BaseRepository[Team]
}

func NewTeamRepository(db *gorm.DB) TeamRepository {
	return &teamRepositoryImpl{
		BaseRepository: utils.NewBaseRepository[Team](db),
	}
}
