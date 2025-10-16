package zone

import (
	"github.com/baimhons/stadiumhub/internal/utils"
	"gorm.io/gorm"
)

type ZoneRepository interface {
	utils.BaseRepository[Zone]
	GetZonesByTeamID(teamID int) ([]Zone, error)
	DB() *gorm.DB
}

type zoneRepositoryImpl struct {
	utils.BaseRepository[Zone]
	db *gorm.DB
}

func NewZoneRepository(db *gorm.DB) ZoneRepository {
	return &zoneRepositoryImpl{
		BaseRepository: utils.NewBaseRepository[Zone](db),
		db:             db,
	}
}

func (zr *zoneRepositoryImpl) GetZonesByTeamID(teamID int) ([]Zone, error) {
	var zones []Zone
	if err := zr.db.Preload("Team").Where("team_id = ?", teamID).Find(&zones).Error; err != nil {
		return nil, err
	}
	return zones, nil
}

func (zr *zoneRepositoryImpl) DB() *gorm.DB {
	return zr.db
}
