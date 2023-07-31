package collectorsrv

import (
	"context"
	"fmt"
	"sync"

	"github.com/mikeewhite/ship-locator/backend/internal/core/domain"
	"github.com/mikeewhite/ship-locator/backend/internal/core/ports"
	"github.com/mikeewhite/ship-locator/backend/pkg/clog"
)

const (
	dataChanSize   = 500
	workerPoolSize = 5
)

type Service struct {
	ch           chan *domain.Ship
	msgPublisher ports.Publisher
	wg           sync.WaitGroup
}

func New(ctx context.Context, publisher ports.Publisher) *Service {
	s := &Service{
		ch:           make(chan *domain.Ship, dataChanSize),
		msgPublisher: publisher,
	}

	for i := 0; i < workerPoolSize; i++ {
		s.wg.Add(1)
		go s.worker(ctx)
	}

	return s
}

func (s *Service) Process(mmsi int32, shipName string, latitude, longitude float64) error {
	ship := domain.NewShip(mmsi, shipName, latitude, longitude)
	if err := ship.Validate(); err != nil {
		return fmt.Errorf("invalid ship entity: %w", err)
	}

	s.ch <- ship

	return nil
}

func (s *Service) Shutdown() {
	s.wg.Wait()
	close(s.ch)
}

func (s *Service) worker(ctx context.Context) {
	defer s.wg.Done()
	for {
		select {
		case ship, ok := <-s.ch:
			if !ok {
				return
			}
			if err := s.msgPublisher.Publish(ship); err != nil {
				clog.Errorw("failed to publish ship",
					"error", err.Error(),
					"mmsi", ship.MMSI,
					"name", ship.Name)
			}
		case <-ctx.Done():
			clog.Info("worker routine stopped")
			return
		}
	}
}
