package graphql

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/mikeewhite/ship-locator/backend/internal/core/domain"
	"github.com/mikeewhite/ship-locator/backend/pkg/apperrors"
	"github.com/mikeewhite/ship-locator/backend/pkg/config"
)

type MockRepo struct {
}

func (mr *MockRepo) Get(_ context.Context, mmsi int32) (domain.Ship, error) {
	if mmsi == 259000420 {
		return domain.Ship{
			MMSI:      259000420,
			Name:      "AUGUSTSON",
			Latitude:  66.02695,
			Longitude: 12.253821666666665,
		}, nil
	}
	return domain.Ship{}, apperrors.NewNoShipFoundErr(mmsi)
}

func (mr *MockRepo) Store(_ context.Context, _ []domain.Ship) error {
	// noop
	return nil
}

func TestHandleQuery_NoMatch(t *testing.T) {
	cfg, err := config.Load()
	require.NoError(t, err)

	srv, err := New(*cfg, &MockRepo{})
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

func TestHandleQuery_FullyPopulatedResponse(t *testing.T) {
	cfg, err := config.Load()
	require.NoError(t, err)

	srv, err := New(*cfg, &MockRepo{})
	require.NoError(t, err)
	defer srv.Shutdown()

	url := `http://localhost:8085/graphql`
	body := `{
			"query": "{ ship(mmsi: 259000420) { mmsi name latitude longitude } }"
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
				"longitude": 12.253821666666665
			}
		}
	}`
	assert.JSONEq(t, expResp, rec.Body.String())
}

func TestHandleQuery_PartiallyPopulatedResponse(t *testing.T) {
	cfg, err := config.Load()
	require.NoError(t, err)

	srv, err := New(*cfg, &MockRepo{})
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
