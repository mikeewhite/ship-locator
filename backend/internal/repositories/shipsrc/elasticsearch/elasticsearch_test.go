package elasticsearch

import (
	"context"
	"github.com/mikeewhite/ship-locator/backend/internal/core/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestSearch_NoMatchingResults(t *testing.T) {
	tv := setup(t)
	matches, err := tv.elasticsearch.Search(context.Background(), "AUGUSTSON")
	require.NoError(t, err)
	assert.Empty(t, matches)
}

func TestSearch_MatchesOnNameAndMMSI(t *testing.T) {
	tv := setup(t)

	// index a single document
	ship := domain.NewShipSearchResult(12345, "AUGUSTSON")
	require.NoError(t, tv.elasticsearch.Index(context.Background(), []domain.ShipSearchResult{ship}))

	// allow time for indexing
	time.Sleep(1 * time.Second)

	tt := map[string]struct {
		query string
	}{
		"exact match on name":            {query: "AUGUSTSON"},
		"case insensitive match on name": {query: "auguston"},
		"partial match on name":          {query: "AUGUST"},
		"fuzzy matches on name":          {query: "AUGUSTEN"},
		"exact match on MMSI":            {query: "12345"},
		"fuzzy match on MMSI":            {query: "12346"},
	}

	for name, tc := range tt {
		t.Run(name, func(t *testing.T) {
			matches, err := tv.elasticsearch.Search(context.Background(), tc.query)
			require.NoError(t, err)
			require.Len(t, matches, 1)
			assert.Equal(t, ship, matches[0])
		})
	}
}

func TestSearch_ReIndexesEntriesWithTheSameMMSI(t *testing.T) {
	tv := setup(t)

	// index a single document
	ship := domain.NewShipSearchResult(259000420, "AUGUSTSON")
	require.NoError(t, tv.elasticsearch.Index(context.Background(), []domain.ShipSearchResult{ship}))

	// allow time for indexing
	time.Sleep(1 * time.Second)

	// re-index the document (to update its name field)
	ship.Name = "AUGUSTSEN"
	require.NoError(t, tv.elasticsearch.Index(context.Background(), []domain.ShipSearchResult{ship}))

	// allow time for indexing
	time.Sleep(1 * time.Second)

	matches, err := tv.elasticsearch.Search(context.Background(), "259000420")
	require.NoError(t, err)
	require.Len(t, matches, 1)
	assert.Equal(t, "AUGUSTSEN", ship.Name)
}
