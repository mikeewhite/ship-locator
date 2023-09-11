package domain

import (
	"errors"
	"strings"
	"time"
)

type Ship struct {
	MMSI        int32
	Name        string
	Latitude    float64
	Longitude   float64
	LastUpdated time.Time
}

func NewShip(mmsi int32, name string, latitude, longitude float64, lastUpdated time.Time) *Ship {
	ship := Ship{
		MMSI:        mmsi,
		Name:        strings.TrimSpace(name),
		Latitude:    latitude,
		Longitude:   longitude,
		LastUpdated: lastUpdated.UTC(),
	}
	return &ship
}

func (s *Ship) Validate() error {
	if s.MMSI == 0 {
		return errors.New("mmsi must be non-zero")
	}

	if s.Latitude < -90 || s.Latitude > 90 {
		return errors.New("invalid latitude")
	}

	if s.Longitude < -180 || s.Longitude > 180 {
		return errors.New("invalid longitude")
	}

	return nil
}
