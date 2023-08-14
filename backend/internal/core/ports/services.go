package ports

import (
	"context"

	"github.com/mikeewhite/ship-locator/backend/internal/core/domain"
)

type CollectorService interface {
	Process(mmsi int32, shipName string, latitude, longitude float64) error
}

type ShipService interface {
	Get(ctx context.Context, mmsi int32) (domain.Ship, error)
	Store(ctx context.Context, ships []domain.Ship) error
}
