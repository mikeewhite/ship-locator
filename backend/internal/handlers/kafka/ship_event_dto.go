package kafka

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/segmentio/kafka-go"

	"github.com/mikeewhite/ship-locator/backend/internal/core/domain"
)

type ShipLocationUpdatedEventDTO struct {
	Key       string  `json:"-"`
	Name      string  `json:"name"`
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

func NewShipLocationUpdatedEventDTOFromDomainEntity(s domain.Ship) *ShipLocationUpdatedEventDTO {
	return &ShipLocationUpdatedEventDTO{
		Key:       strconv.FormatInt(int64(s.MMSI), 10),
		Name:      s.Name,
		Latitude:  s.Latitude,
		Longitude: s.Longitude,
	}
}

func NewShipLocationUpdatedEventDTOFromKafkaMsg(msg *kafka.Message) (*ShipLocationUpdatedEventDTO, error) {
	var dto ShipLocationUpdatedEventDTO
	err := json.Unmarshal(msg.Value, &dto)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal ship data: %w", err)
	}
	dto.Key = string(msg.Key)
	return &dto, err
}

func (dto *ShipLocationUpdatedEventDTO) ToDomainEntity() (*domain.ShipSearchResult, error) {
	mmsi, err := strconv.ParseInt(dto.Key, 10, 32)
	if err != nil {
		return nil, fmt.Errorf("failed to convert key '%s' to integer: %w", dto.Key, err)
	}
	ssr := domain.NewShipSearchResult(int32(mmsi), dto.Name)
	return &ssr, nil
}
