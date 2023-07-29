package collectorsrv

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/mikeewhite/ship-locator/backend/internal/core/domain"
)

type MockProducer struct {
	queue []*domain.Ship
}

func (mp *MockProducer) Publish(data *domain.Ship) error {
	mp.queue = append(mp.queue, data)
	return nil
}

func TestService_Process(t *testing.T) {
	mockProducer := &MockProducer{}
	s := New(context.Background(), mockProducer)

	require.NoError(t, s.Process(12345, "CALL SIGN", 66.02695, 12.253821666666665))

	time.Sleep(200) // sleep for 200 ms to allow time for worker pool to process job

	require.NotEmpty(t, mockProducer.queue)
	require.Len(t, mockProducer.queue, 1)
	ship := mockProducer.queue[0]
	assert.Equal(t, int32(12345), ship.MMSI)
	assert.Equal(t, "CALL SIGN", ship.Name)
	assert.Equal(t, 66.02695, ship.Latitude)
	assert.Equal(t, 12.253821666666665, ship.Longitude)
}
