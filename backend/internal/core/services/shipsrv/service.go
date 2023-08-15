package shipsrv

import (
	"context"

	"github.com/mikeewhite/ship-locator/backend/internal/core/domain"
	"github.com/mikeewhite/ship-locator/backend/internal/core/ports"
)

type Service struct {
	repo ports.ShipRepository
}

func New(repo ports.ShipRepository) *Service {
	return &Service{repo: repo}
}

func (s *Service) Get(ctx context.Context, mmsi int32) (domain.Ship, error) {
	return s.repo.Get(ctx, mmsi)
}

func (s *Service) Store(ctx context.Context, ships []domain.Ship) error {
	return s.repo.Store(ctx, ships)
}
