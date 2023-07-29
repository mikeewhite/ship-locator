package domain

import "errors"

type Ship struct {
	MMSI      int32
	Name      string
	Latitude  float64
	Longitude float64
}

func NewShip(mmsi int32, name string, latitude, longitude float64) *Ship {
	ship := Ship{
		MMSI:      mmsi,
		Name:      name,
		Latitude:  latitude,
		Longitude: longitude,
	}
	return &ship
}

func (s *Ship) Validate() error {
	if s.MMSI == 0 {
		return errors.New("mmsi must be non-zero")
	}
	return nil
}
