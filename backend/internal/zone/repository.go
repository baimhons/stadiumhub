package zone

import "github.com/baimhons/stadiumhub/internal/utils"

type ZoneRepository interface {
	utils.BaseRepository[Zone]
}

type zoneRepositoryImpl struct {
	utils.BaseRepository[Zone]
}

func NewZoneRepository() ZoneRepository {
	return &zoneRepositoryImpl{
		BaseRepository: utils.NewBaseRepository[Zone](nil),
	}
}
