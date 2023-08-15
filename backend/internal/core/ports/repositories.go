package ports

import (
	"context"

	"github.com/mikeewhite/ship-locator/backend/internal/core/domain"
)

type ShipRepository interface {
	Get(ctx context.Context, mmsi int32) (domain.Ship, error)
	Store(ctx context.Context, ships []domain.Ship) error
}
