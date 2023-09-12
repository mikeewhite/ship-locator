package kafka

import (
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/segmentio/kafka-go"

	"github.com/mikeewhite/ship-locator/backend/internal/core/domain"
)

type shipDTO struct {
	Key       string  `json:"-"`
	Name      string  `json:"name"`
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

func newShipDTOFromDomainEntity(s domain.Ship) *shipDTO {
	return &shipDTO{
		Key:       strconv.FormatInt(int64(s.MMSI), 10),
		Name:      s.Name,
		Latitude:  s.Latitude,
		Longitude: s.Longitude,
	}
}

func newShipDTOFromKafkaMsg(msg *kafka.Message) (*shipDTO, error) {
	var dto shipDTO
	err := json.Unmarshal(msg.Value, &dto)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal ship data: %w", err)
	}
	dto.Key = string(msg.Key)
	return &dto, err
}

func (dto *shipDTO) toDomainEntity() (*domain.Ship, error) {
	mmsi, err := strconv.ParseInt(dto.Key, 10, 32)
	if err != nil {
		return nil, fmt.Errorf("failed to convert key '%s' to integer: %w", dto.Key, err)
	}
	return domain.NewShip(int32(mmsi), dto.Name, dto.Latitude, dto.Longitude, time.Now()), nil
}
