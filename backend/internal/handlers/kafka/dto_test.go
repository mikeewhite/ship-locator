package kafka

import (
	"testing"

	"github.com/segmentio/kafka-go"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/mikeewhite/ship-locator/backend/internal/core/domain"
)

func TestNewShipDTOFromDomainEntity(t *testing.T) {
	s := domain.Ship{
		MMSI:      259000420,
		Name:      "AUGUSTSON",
		Latitude:  66.02695,
		Longitude: 12.253821666666665,
	}

	dto := newShipDTOFromDomainEntity(s)
	assert.Equal(t, "259000420", dto.Key)
	assert.Equal(t, "AUGUSTSON", dto.Name)
	assert.Equal(t, 66.02695, dto.Latitude)
	assert.Equal(t, 12.253821666666665, dto.Longitude)
}

func TestNewShipDTOFromKafkaMsg(t *testing.T) {
	msg := &kafka.Message{
		Key:   []byte("259000420"),
		Value: []byte(`{"name":"AUGUSTSON", "latitude":66.02695, "longitude":12.253821666666665}`),
	}

	dto, err := newShipDTOFromKafkaMsg(msg)
	require.NoError(t, err)
	assert.Equal(t, "259000420", dto.Key)
	assert.Equal(t, "AUGUSTSON", dto.Name)
	assert.Equal(t, 66.02695, dto.Latitude)
	assert.Equal(t, 12.253821666666665, dto.Longitude)
}

func TestToDomainEntity(t *testing.T) {
	dto := &shipDTO{
		Key:       "259000420",
		Name:      "AUGUSTSON",
		Latitude:  66.02695,
		Longitude: 12.253821666666665,
	}

	entity, err := dto.toDomainEntity()
	require.NoError(t, err)
	assert.Equal(t, int32(259000420), entity.MMSI)
	assert.Equal(t, "AUGUSTSON", entity.Name)
	assert.Equal(t, 66.02695, entity.Latitude)
	assert.Equal(t, 12.253821666666665, entity.Longitude)
}
