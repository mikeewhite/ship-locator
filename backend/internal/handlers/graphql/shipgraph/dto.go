package shipgraph

import (
	"time"

	"github.com/mikeewhite/ship-locator/backend/internal/core/domain"
)

type Ship struct {
	MMSI        int32     `json:"mmsi"`
	Name        string    `json:"name"`
	Latitude    float64   `json:"latitude"`
	Longitude   float64   `json:"longitude"`
	LastUpdated time.Time `json:"lastUpdated"`
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
