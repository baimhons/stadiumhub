package zone

import (
	"errors"
	"net/http"

	"github.com/baimhons/stadiumhub/internal/match"
	"github.com/baimhons/stadiumhub/internal/zone/api/response"
	"gorm.io/gorm"
)

type ZoneService interface {
	GetZoneByMatchID(matchID int) (resp []*response.ZoneResponse, statusCode int, err error)
}

type zoneServiceImpl struct {
	zoneRepository ZoneRepository
}

func NewZoneService(zoneRepository ZoneRepository) ZoneService {
	return &zoneServiceImpl{
		zoneRepository: zoneRepository,
	}
}

func (zs *zoneServiceImpl) GetZoneByMatchID(matchID int) (resp []*response.ZoneResponse, statusCode int, err error) {
	var m match.Match
	if err := zs.zoneRepository.DB().First(&m, matchID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, http.StatusBadRequest, errors.New("match not found")
		}
		return nil, http.StatusInternalServerError, err
	}

	zones, err := zs.zoneRepository.GetZonesByTeamID(m.HomeTeamID)
	if err != nil {
		return nil, http.StatusInternalServerError, err
	}

	for _, z := range zones {
		resp = append(resp, &response.ZoneResponse{
			ID:   z.ID,
			Team: z.Team,
			Name: z.Name,
		})
	}

	return resp, http.StatusOK, nil
}
