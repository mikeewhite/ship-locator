package consumer

import (
	"context"
	"fmt"
	"time"

	"github.com/segmentio/kafka-go"

	"github.com/mikeewhite/ship-locator/backend/internal/core/domain"
	"github.com/mikeewhite/ship-locator/backend/internal/core/ports"
	kafka2 "github.com/mikeewhite/ship-locator/backend/internal/handlers/kafka"
	"github.com/mikeewhite/ship-locator/backend/pkg/clog"
	"github.com/mikeewhite/ship-locator/backend/pkg/config"
)

type ShipEventConsumer struct {
	reader  *kafka.Reader
	service ports.ShipSearchService
	metrics kafka2.Metrics
}

func NewShipEventConsumer(cfg config.Config, service ports.ShipSearchService, metrics kafka2.Metrics) (*ShipEventConsumer, error) {
	return &ShipEventConsumer{
		reader: kafka.NewReader(kafka.ReaderConfig{
			Brokers:  []string{cfg.KafkaAddress},
			GroupID:  cfg.KafkaConsumerGroup,
			Topic:    cfg.KafkaShipEventTopic,
			MaxBytes: 10e6, // 10MB
			ErrorLogger: kafka.LoggerFunc(func(msg string, a ...interface{}) {
				clog.Errorf(msg, a...)
				fmt.Println()
			}),
		}),
		service: service,
		metrics: metrics,
	}, nil
}

func (c *ShipEventConsumer) Read(ctx context.Context) error {
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

			dto, err := kafka2.NewShipLocationUpdatedEventDTOFromKafkaMsg(&m)
			if err != nil {
				return fmt.Errorf("error on generating DTO from Kafka message: %w", err)
			}
			clog.Infof("🚢: %v", dto)

			shipSearchResult, err := dto.ToDomainEntity()
			if err != nil {
				return fmt.Errorf("error on converting ship DTO to domain entity: %w", err)
			}

			err = c.service.Store(ctx, []domain.ShipSearchResult{*shipSearchResult})
			if err != nil {
				return fmt.Errorf("error on storing ship search result: %w", err)
			}
		}
	}
}

func (c *ShipEventConsumer) Shutdown() {
	if err := c.reader.Close(); err != nil {
		clog.Errorf("failed to close Kafka reader: %s", err.Error())
	}
}
