package kafka

import (
	"context"
	"fmt"

	"github.com/segmentio/kafka-go"

	"github.com/mikeewhite/ship-locator/backend/internal/core/domain"
	"github.com/mikeewhite/ship-locator/backend/internal/core/ports"
	"github.com/mikeewhite/ship-locator/backend/pkg/clog"
	"github.com/mikeewhite/ship-locator/backend/pkg/config"
)

type Consumer struct {
	reader  *kafka.Reader
	service ports.ShipService
}

func NewConsumer(cfg config.Config, service ports.ShipService) (*Consumer, error) {
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
		service: service,
	}, nil
}

func (c *Consumer) Read(ctx context.Context) error {
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			m, err := c.reader.ReadMessage(ctx)
			if err != nil {
				return fmt.Errorf("error on reading message: %w", err)
			}

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
		}
	}
}

func (c *Consumer) Shutdown() {
	if err := c.reader.Close(); err != nil {
		clog.Errorf("failed to close Kafka reader: %s", err.Error())
	}
}
