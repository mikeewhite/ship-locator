package elasticsearch

import (
	"encoding/json"
	"github.com/mikeewhite/ship-locator/backend/internal/core/domain"
)

type shipDTO struct {
	MMSI int32  `json:"mmsi"`
	Name string `json:"name"`
}

func toShipDTO(s domain.ShipSearchResult) shipDTO {
	return shipDTO{
		MMSI: s.MMSI,
		Name: s.Name,
	}
}

func (s *shipDTO) toJSON() ([]byte, error) {
	return json.Marshal(s)
}
