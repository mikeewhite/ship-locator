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

type ShipSearchService interface {
	Search(ctx context.Context, searchTerm string) ([]domain.ShipSearchResult, error)
	Store(ctx context.Context, ships []domain.ShipSearchResult) error
}
