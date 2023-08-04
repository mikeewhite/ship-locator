package kafka

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/segmentio/kafka-go"

	"github.com/mikeewhite/ship-locator/backend/internal/core/domain"
	"github.com/mikeewhite/ship-locator/backend/pkg/clog"
	"github.com/mikeewhite/ship-locator/backend/pkg/config"
)

type Producer struct {
	writer *kafka.Writer
}

func NewProducer(cfg config.Config) (*Producer, error) {
	return &Producer{
		writer: &kafka.Writer{
			Addr:     kafka.TCP(cfg.KafkaAddress),
			Topic:    cfg.KafkaTopic,
			Balancer: &kafka.LeastBytes{},
			ErrorLogger: kafka.LoggerFunc(func(msg string, a ...interface{}) {
				clog.Errorf(msg, a...)
				fmt.Println()
			}),
		},
	}, nil
}

func (p *Producer) Write(ctx context.Context, data domain.Ship) error {
	dto := newShipDTOFromDomainEntity(data)
	b, err := json.Marshal(dto)
	if err != nil {
		return fmt.Errorf("failed to marshal ship DTO: %w", err)
	}
	return p.writer.WriteMessages(ctx, kafka.Message{Key: []byte(dto.Key), Value: b})
}

func (p *Producer) Shutdown() {
	if err := p.writer.Close(); err != nil {
		clog.Errorf("failed to close Kafka writer: %s", err.Error())
	}
}
