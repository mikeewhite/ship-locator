package postgres

import (
	"context"
	"testing"
	"time"

	"github.com/mikeewhite/ship-locator/backend/pkg/config"
)

type testVars struct {
	pg *Postgres
}

type NoopMetricsClient struct {
}

func (mc *NoopMetricsClient) DBQueryTime(queryName string, startTime time.Time) {}

func setup(t *testing.T) *testVars {
	t.Helper()

	cfg, err := config.Load()
	if err != nil {
		t.Fatal(err)
	}
	cfg.PostgresAddress = "localhost:5433"

	pg, err := NewPostgres(context.Background(), *cfg, &NoopMetricsClient{})
	if err != nil {
		t.Fatal(err)
	}

	t.Cleanup(func() {
		// wipe database
		_, err := pg.conn.Exec(context.Background(), "DELETE FROM ships")
		if err != nil {
			t.Fatal(err)
		}

		_ = pg.conn.Close(context.Background())
	})
	return &testVars{pg: pg}
}
