package shipgraph

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/mikeewhite/ship-locator/backend/internal/core/domain"
	"github.com/mikeewhite/ship-locator/backend/pkg/apperrors"
	"github.com/mikeewhite/ship-locator/backend/pkg/config"
)

type MockShipService struct {
}

func (mr *MockShipService) Get(_ context.Context, mmsi int32) (domain.Ship, error) {
	lastUpdated, _ := time.Parse(time.RFC3339, "2023-09-11T17:04:05Z")
	if mmsi == 259000420 {
		return domain.Ship{
			MMSI:        259000420,
			Name:        "AUGUSTSON",
			Latitude:    66.02695,
			Longitude:   12.253821666666665,
			LastUpdated: lastUpdated,
		}, nil
	}
	return domain.Ship{}, apperrors.NewNoShipFoundErr(mmsi)
}

func (mr *MockShipService) Store(_ context.Context, _ []domain.Ship) error {
	// noop
	return nil
}

func TestHandleQuery_Ship_NoMatch(t *testing.T) {
	cfg, err := config.Load()
	require.NoError(t, err)

	srv, err := New(*cfg, &MockShipService{})
	require.NoError(t, err)
	defer srv.Shutdown()

	url := `http://localhost:8085/graphql`
	body := `{
			"query": "{ ship(mmsi: 1234) { mmsi name latitude longitude } }"
		}`
	req, err := http.NewRequest("POST", url, strings.NewReader(body))
	require.NoError(t, err)

	rec := httptest.NewRecorder()
	srv.HandleQuery(rec, req)

	assert.Equal(t, 200, rec.Code)
	assert.Contains(t, rec.Body.String(), "no matching ship found for id: 1234")
}

func TestHandleQuery_Ship_FullyPopulatedResponse(t *testing.T) {
	cfg, err := config.Load()
	require.NoError(t, err)

	srv, err := New(*cfg, &MockShipService{})
	require.NoError(t, err)
	defer srv.Shutdown()

	url := `http://localhost:8085/graphql`
	body := `{
			"query": "{ ship(mmsi: 259000420) { mmsi name latitude longitude lastUpdated } }"
		}`
	req, err := http.NewRequest("POST", url, strings.NewReader(body))
	req.Header.Add("Content-Type", "application/json")
	require.NoError(t, err)

	rec := httptest.NewRecorder()
	srv.HandleQuery(rec, req)

	require.Equal(t, 200, rec.Code)
	expResp := `{
		"data": {
			"ship": {
				"mmsi": 259000420,
				"name": "AUGUSTSON",
				"latitude": 66.02695,
				"longitude": 12.253821666666665,
				"lastUpdated": "2023-09-11T17:04:05Z"
			}
		}
	}`
	assert.JSONEq(t, expResp, rec.Body.String())
}

func TestHandleQuery_Ship_PartiallyPopulatedResponse(t *testing.T) {
	cfg, err := config.Load()
	require.NoError(t, err)

	srv, err := New(*cfg, &MockShipService{})
	require.NoError(t, err)
	defer srv.Shutdown()

	url := `http://localhost:8085/graphql`
	body := `{
			"query": "{ ship(mmsi: 259000420) { name } }"
		}`
	req, err := http.NewRequest("POST", url, strings.NewReader(body))
	require.NoError(t, err)

	rec := httptest.NewRecorder()
	srv.HandleQuery(rec, req)

	assert.Equal(t, 200, rec.Code)
	expResp := `{
		"data": {
			"ship": {
				"name": "AUGUSTSON"
			}
		}
	}`
	assert.JSONEq(t, expResp, rec.Body.String())
}
