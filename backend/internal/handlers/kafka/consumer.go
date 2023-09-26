package kafka

import (
	"context"
	"fmt"
	"time"

	"github.com/segmentio/kafka-go"

	"github.com/mikeewhite/ship-locator/backend/internal/core/domain"
	"github.com/mikeewhite/ship-locator/backend/internal/core/ports"
	"github.com/mikeewhite/ship-locator/backend/pkg/clog"
	"github.com/mikeewhite/ship-locator/backend/pkg/config"
)

type Consumer struct {
	reader        *kafka.Reader
	service       ports.ShipService
	searchService ports.ShipSearchService
	metrics       Metrics
}

func NewConsumer(cfg config.Config, service ports.ShipService, searchService ports.ShipSearchService, metrics Metrics) (*Consumer, error) {
	return &Consumer{
		reader: kafka.NewReader(kafka.ReaderConfig{
			Brokers:  []string{cfg.KafkaAddress},
			GroupID:  cfg.KafkaConsumerGroup,
			Topic:    cfg.KafkaTopic,
			MaxBytes: 10e6, // 10MB
			ErrorLogger: kafka.LoggerFunc(func(msg string, a ...interface{}) {
				clog.Errorf(msg, a...)
				fmt.Println()
			}),
		}),
		service:       service,
		searchService: searchService,
		metrics:       metrics,
	}, nil
}

func (c *Consumer) Read(ctx context.Context) error {
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			start := time.Now()
			m, err := c.reader.ReadMessage(ctx)
			if err != nil {
				return fmt.Errorf("error on reading message: %w", err)
			}
			c.metrics.KafkaConsumeTime(c.reader.Config().Topic, start)

			dto, err := newShipDTOFromKafkaMsg(&m)
			if err != nil {
				return fmt.Errorf("error on generating DTO from Kafka message: %w", err)
			}
			clog.Infof("🚢: %v", dto)

			ship, err := dto.toDomainEntity()
			if err != nil {
				return fmt.Errorf("error on converting ship DTO to domain entity: %w", err)
			}
			ships := []domain.Ship{*ship}
			err = c.service.Store(ctx, ships)
			if err != nil {
				return fmt.Errorf("error on storing ship data: %w", err)
			}

			// TODO - this should be a separate consumer for the search service which should be consuming 'ship-data-stored' events on a different topic
			searchResult, err := dto.toDomainSearchResult()
			if err != nil {
				return fmt.Errorf("error on converting ship DTO to domain search result: %w", err)
			}
			err = c.searchService.Store(ctx, []domain.ShipSearchResult{*searchResult})
			if err != nil {
				return fmt.Errorf("error on storing ship search result: %w", err)
			}
		}
	}
}

func (c *Consumer) Shutdown() {
	if err := c.reader.Close(); err != nil {
		clog.Errorf("failed to close Kafka reader: %s", err.Error())
	}
}
