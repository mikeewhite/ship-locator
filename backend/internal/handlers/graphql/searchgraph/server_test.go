package searchgraph

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/mikeewhite/ship-locator/backend/internal/core/domain"
	"github.com/mikeewhite/ship-locator/backend/pkg/config"
)

type MockShipSearchService struct {
}

func (msss *MockShipSearchService) Search(_ context.Context, searchTerm string) ([]domain.ShipSearchResult, error) {
	if searchTerm == "AUGUSTSON" {
		return []domain.ShipSearchResult{
			{
				MMSI: 259000420,
				Name: "AUGUSTSON",
			},
		}, nil
	}
	return []domain.ShipSearchResult{}, nil
}

func (msss *MockShipSearchService) Store(_ context.Context, _ []domain.ShipSearchResult) error {
	// noop
	return nil
}

func TestHandleQuery_ShipSearch_NoMatch(t *testing.T) {
	cfg, err := config.Load()
	require.NoError(t, err)

	srv, err := New(*cfg, &MockShipSearchService{})
	require.NoError(t, err)
	defer srv.Shutdown()

	url := `http://localhost:8085/graphql`
	body := `{
			"query": "{ shipSearch(searchTerm: \"1234\") { mmsi name } }"
		}`
	req, err := http.NewRequest("POST", url, strings.NewReader(body))
	require.NoError(t, err)

	rec := httptest.NewRecorder()
	srv.HandleQuery(rec, req)

	assert.Equal(t, 200, rec.Code)
	assert.Contains(t, rec.Body.String(), "")
}

func TestHandleQuery_ShipSearch_FullyPopulatedResponse(t *testing.T) {
	cfg, err := config.Load()
	require.NoError(t, err)

	srv, err := New(*cfg, &MockShipSearchService{})
	require.NoError(t, err)
	defer srv.Shutdown()

	url := `http://localhost:8085/graphql`
	body := `{
			"query": "{ shipSearch(searchTerm: \"AUGUSTSON\") { mmsi name } }"
		}`
	req, err := http.NewRequest("POST", url, strings.NewReader(body))
	req.Header.Add("Content-Type", "application/json")
	require.NoError(t, err)

	rec := httptest.NewRecorder()
	srv.HandleQuery(rec, req)

	require.Equal(t, 200, rec.Code)
	expResp := `{
		"data": {
			"shipSearch": [{
				"mmsi": 259000420,
				"name": "AUGUSTSON"
			}]
		}
	}`
	assert.JSONEq(t, expResp, rec.Body.String())
}

func TestHandleQuery_ShipSearch_PartiallyPopulatedResponse(t *testing.T) {
	cfg, err := config.Load()
	require.NoError(t, err)

	srv, err := New(*cfg, &MockShipSearchService{})
	require.NoError(t, err)
	defer srv.Shutdown()

	url := `http://localhost:8085/graphql`
	body := `{
			"query": "{ shipSearch(searchTerm: \"AUGUSTSON\") { name } }"
		}`
	req, err := http.NewRequest("POST", url, strings.NewReader(body))
	require.NoError(t, err)

	rec := httptest.NewRecorder()
	srv.HandleQuery(rec, req)

	assert.Equal(t, 200, rec.Code)
	expResp := `{
		"data": {
			"shipSearch": [{
				"name": "AUGUSTSON"
			}]
		}
	}`
	assert.JSONEq(t, expResp, rec.Body.String())
}
