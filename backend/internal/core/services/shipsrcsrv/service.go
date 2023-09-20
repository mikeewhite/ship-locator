package shipsrcsrv

import (
	"context"
	"github.com/mikeewhite/ship-locator/backend/internal/core/domain"
	"github.com/mikeewhite/ship-locator/backend/internal/core/ports"
)

type Service struct {
	repo ports.ShipSearchRepository
}

func New(repo ports.ShipSearchRepository) *Service {
	return &Service{repo: repo}
}

func (s *Service) Search(ctx context.Context, searchTerm string) ([]domain.ShipSearchResult, error) {
	return s.repo.Search(ctx, searchTerm)
}

func (s *Service) Store(ctx context.Context, ships []domain.ShipSearchResult) error {
	return s.repo.Index(ctx, ships)
}
