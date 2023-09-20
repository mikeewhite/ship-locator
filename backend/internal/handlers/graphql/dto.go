package graphql

import (
	"github.com/mikeewhite/ship-locator/backend/internal/core/domain"
	"time"
)

type Ship struct {
	MMSI        int32     `json:"mmsi"`
	Name        string    `json:"name"`
	Latitude    float64   `json:"latitude"`
	Longitude   float64   `json:"longitude"`
	LastUpdated time.Time `json:"lastUpdated"`
}

type ShipSearchResult struct {
	MMSI int32  `json:"mmsi"`
	Name string `json:"name"`
}

func toDTO(s domain.Ship) Ship {
	return Ship{
		MMSI:        s.MMSI,
		Name:        s.Name,
		Latitude:    s.Latitude,
		Longitude:   s.Longitude,
		LastUpdated: s.LastUpdated,
	}
}

func toShipResultDTOs(searchResults []domain.ShipSearchResult) []ShipSearchResult {
	dtos := make([]ShipSearchResult, len(searchResults))
	for i := 0; i < len(searchResults); i++ {
		dtos[i] = ShipSearchResult{
			MMSI: searchResults[i].MMSI,
			Name: searchResults[i].Name,
		}
	}

	return dtos
}
