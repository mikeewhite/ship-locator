package producer

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/segmentio/kafka-go"

	"github.com/mikeewhite/ship-locator/backend/internal/core/domain"
	kafka2 "github.com/mikeewhite/ship-locator/backend/internal/handlers/kafka"
	"github.com/mikeewhite/ship-locator/backend/pkg/clog"
	"github.com/mikeewhite/ship-locator/backend/pkg/config"
)

type ShipEventProducer struct {
	writer *kafka.Writer
}

func NewShipEventProducer(cfg config.Config) (*ShipEventProducer, error) {
	return &ShipEventProducer{
		writer: &kafka.Writer{
			Addr:     kafka.TCP(cfg.KafkaAddress),
			Topic:    cfg.KafkaShipEventTopic,
			Balancer: &kafka.LeastBytes{},
			ErrorLogger: kafka.LoggerFunc(func(msg string, a ...interface{}) {
				clog.Errorf(msg, a...)
				fmt.Println()
			}),
		},
	}, nil
}

func (p *ShipEventProducer) PublishShipLocationsUpdatedEvent(ctx context.Context, ships []domain.Ship) error {
	// For simplicity, we'll publish a single message per event rather than a bulk message using the
	// MMSI of the ship as the key. This will ensure that all events for a given ship are processed
	// in order by a single consumer. It will also allow us to potentially use a compacted topic so to only
	// keep the latest message for a given ship.
	for _, ship := range ships {
		event := kafka2.NewShipLocationUpdatedEventDTOFromDomainEntity(ship)
		b, err := json.Marshal(event)
		if err != nil {
			return fmt.Errorf("failed to marshal ship event DTO: %w", err)
		}
		if err := p.writer.WriteMessages(ctx, kafka.Message{Key: []byte(event.Key), Value: b}); err != nil {
			return fmt.Errorf("failed to publish ship location updated event: %w", err)
		}
	}

	return nil
}

func (p *ShipEventProducer) Shutdown() {
	if err := p.writer.Close(); err != nil {
		clog.Errorf("failed to close Kafka writer: %s", err.Error())
	}
}
