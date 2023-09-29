package shipsrv

import (
	"context"
	"fmt"

	"github.com/mikeewhite/ship-locator/backend/internal/core/domain"
	"github.com/mikeewhite/ship-locator/backend/internal/core/ports"
)

type ShipEventProducer interface {
	PublishShipLocationsUpdatedEvent(ctx context.Context, ships []domain.Ship) error
}

type Service struct {
	repo              ports.ShipRepository
	shipEventProducer ShipEventProducer
}

func New(repo ports.ShipRepository, shipEventProducer ShipEventProducer) *Service {
	return &Service{
		repo:              repo,
		shipEventProducer: shipEventProducer,
	}
}

func (s *Service) Get(ctx context.Context, mmsi int32) (domain.Ship, error) {
	return s.repo.Get(ctx, mmsi)
}

func (s *Service) Store(ctx context.Context, ships []domain.Ship) error {
	err := s.repo.Store(ctx, ships)
	if err != nil {
		return fmt.Errorf("failed to store ships: %w", err)
	}

	return s.shipEventProducer.PublishShipLocationsUpdatedEvent(ctx, ships)
}
