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

type ShipDataProducer struct {
	writer *kafka.Writer
}

func NewShipDataProducer(cfg config.Config) (*ShipDataProducer, error) {
	return &ShipDataProducer{
		writer: &kafka.Writer{
			Addr:     kafka.TCP(cfg.KafkaAddress),
			Topic:    cfg.KafkaShipDataTopic,
			Balancer: &kafka.LeastBytes{},
			ErrorLogger: kafka.LoggerFunc(func(msg string, a ...interface{}) {
				clog.Errorf(msg, a...)
				fmt.Println()
			}),
		},
	}, nil
}

func (p *ShipDataProducer) Write(ctx context.Context, data domain.Ship) error {
	dto := kafka2.NewShipDTOFromDomainEntity(data)
	b, err := json.Marshal(dto)
	if err != nil {
		return fmt.Errorf("failed to marshal ship DTO: %w", err)
	}
	return p.writer.WriteMessages(ctx, kafka.Message{Key: []byte(dto.Key), Value: b})
}

func (p *ShipDataProducer) Shutdown() {
	if err := p.writer.Close(); err != nil {
		clog.Errorf("failed to close Kafka writer: %s", err.Error())
	}
}
