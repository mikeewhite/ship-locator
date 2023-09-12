package postgres

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/mikeewhite/ship-locator/backend/internal/core/domain"
	"github.com/mikeewhite/ship-locator/backend/pkg/apperrors"
)

func TestGet_NoMatchingResult(t *testing.T) {
	tv := setup(t)
	_, err := tv.pg.Get(context.Background(), 12345)
	assert.NotNil(t, err)
	expErr := apperrors.NewNoShipFoundErr(12345)
	assert.ErrorAs(t, err, &expErr)
}

func TestStore(t *testing.T) {
	timestamp, _ := time.Parse(time.RFC3339, "2023-09-11T17:04:05Z")
	ship := domain.Ship{
		MMSI:        259000420,
		Name:        "AUGUSTSON",
		Latitude:    66.02695,
		Longitude:   12.253821666666665,
		LastUpdated: timestamp,
	}
	ships := []domain.Ship{ship}

	tv := setup(t)
	err := tv.pg.Store(context.Background(), ships)
	require.NoError(t, err)

	returnedShip, err := tv.pg.Get(context.Background(), 259000420)
	require.NoError(t, err)

	assert.Equal(t, int32(259000420), returnedShip.MMSI)
	assert.Equal(t, "AUGUSTSON", returnedShip.Name)
	assert.Equal(t, 66.02695, returnedShip.Latitude)
	assert.Equal(t, 12.253821666666665, returnedShip.Longitude)
	assert.Equal(t, timestamp, returnedShip.LastUpdated)
}

func TestStore_OverwritesExistingEntries(t *testing.T) {
	now := time.Now().UTC()
	ship := domain.Ship{
		MMSI:        259000420,
		Name:        "AUGUSTSON",
		Latitude:    66.02695,
		Longitude:   12.253821666666665,
		LastUpdated: now,
	}
	ships := []domain.Ship{ship}

	tv := setup(t)
	require.NoError(t, tv.pg.Store(context.Background(), ships))

	// change the ship location and store again
	ships[0].Latitude = 66.03421
	ships[0].Longitude = 12.34251
	ships[0].LastUpdated = time.Now().UTC()
	require.NoError(t, tv.pg.Store(context.Background(), ships))

	// retrieve the entry to check its details have been updated
	updatedShip, err := tv.pg.Get(context.Background(), 259000420)
	require.NoError(t, err)
	assert.Equal(t, 66.03421, updatedShip.Latitude)
	assert.Equal(t, 12.34251, updatedShip.Longitude)

	// check that the updated_at field has been updated
	assert.True(t, updatedShip.LastUpdated.After(now))
}
