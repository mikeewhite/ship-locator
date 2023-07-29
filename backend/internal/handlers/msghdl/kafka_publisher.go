package msghdl

import (
	"github.com/mikeewhite/ship-locator/backend/internal/core/domain"
	"github.com/mikeewhite/ship-locator/backend/pkg/clog"
)

type KafkaPublisher struct {
}

// TODO - accept a logger instance to pass to Kafka driver
func NewKafkaPublisher() *KafkaPublisher {
	return &KafkaPublisher{}
}

func (kp *KafkaPublisher) Publish(data *domain.Ship) error {
	clog.Infow("🚢",
		"mmsi", data.MMSI,
		"name", data.Name,
		"longitude", data.Longitude,
		"latitude", data.Latitude)

	// TODO - publish to Kafka (using a DTO)
	// TODO - consider adding to batch and periodically flushing
	return nil
}

func (kp *KafkaPublisher) Shutdown() {
	// TODO
}
