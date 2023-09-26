package elasticsearch

import (
	"context"
	"github.com/mikeewhite/ship-locator/backend/pkg/config"
	"testing"
)

type testVars struct {
	elasticsearch *Repository
}

func setup(t *testing.T) *testVars {
	t.Helper()

	cfg, err := config.Load()
	if err != nil {
		t.Fatal(err)
	}
	cfg.ElasticsearchIndex = "test_ship_search_index"

	elasticsearch, err := New(context.Background(), *cfg)
	if err != nil {
		t.Fatalf("failed to create elasticsearch client: %s", err)
	}

	t.Cleanup(func() {
		_, err := elasticsearch.client.Indices.Delete(elasticsearch.indexName).Do(context.Background())
		if err != nil {
			t.Fatalf("failed to delete index: %s", err)
		}
	})

	return &testVars{
		elasticsearch: elasticsearch,
	}
}
