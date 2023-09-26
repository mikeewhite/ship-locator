package ports

import (
	"context"

	"github.com/mikeewhite/ship-locator/backend/internal/core/domain"
)

type ShipRepository interface {
	Get(ctx context.Context, mmsi int32) (domain.Ship, error)
	Store(ctx context.Context, ships []domain.Ship) error
}

type ShipSearchRepository interface {
	Search(ctx context.Context, query string) ([]domain.ShipSearchResult, error)
	Index(ctx context.Context, ships []domain.ShipSearchResult) error
}
