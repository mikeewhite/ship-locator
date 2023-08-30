package postgres

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/mikeewhite/ship-locator/backend/internal/core/domain"
	"github.com/mikeewhite/ship-locator/backend/pkg/apperrors"
	"github.com/mikeewhite/ship-locator/backend/pkg/clog"
	"github.com/mikeewhite/ship-locator/backend/pkg/config"
)

type Metrics interface {
	DBQueryTime(queryName string, startTime time.Time)
}

type Postgres struct {
	pool    *pgxpool.Pool
	metrics Metrics
}

const (
	selectSQL = `
			SELECT name, latitude, longitude
			FROM ships
			WHERE mmsi=$1`
	updateSQL = `
			INSERT INTO ships (mmsi, name, latitude, longitude)
			VALUES ($1, $2, $3, $4)
			ON CONFLICT (mmsi)
			DO
				UPDATE SET name = EXCLUDED.name, latitude = EXCLUDED.latitude, longitude = EXCLUDED.longitude`
)

func NewPostgres(ctx context.Context, cfg config.Config, metrics Metrics) (*Postgres, error) {
	url := fmt.Sprintf("postgres://%s:%s@%s/%s",
		cfg.PostgresUsername, cfg.PostgresPassword, cfg.PostgresAddress, cfg.PostgresDBName)

	poolCfg, err := pgxpool.ParseConfig(url)
	if err != nil {
		return nil, fmt.Errorf("failed to parse postgres config: %w", err)
	}
	poolCfg.ConnConfig.Tracer = newTracer()
	pool, err := pgxpool.NewWithConfig(ctx, poolCfg)
	if err != nil {
		return nil, fmt.Errorf("failed to establish connection to postgres: %w", err)
	}

	// check the connection is healthy
	err = pool.Ping(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to ping postgres: %w", err)
	}

	return &Postgres{pool: pool, metrics: metrics}, nil
}

func (pg *Postgres) Get(ctx context.Context, mmsi int32) (domain.Ship, error) {
	defer pg.metrics.DBQueryTime("get_ship_data", time.Now())
	var name string
	var latitude float64
	var longitude float64

	err := pg.pool.QueryRow(ctx, selectSQL, mmsi).Scan(&name, &latitude, &longitude)
	if err == pgx.ErrNoRows {
		return domain.Ship{}, apperrors.NewNoShipFoundErr(mmsi)
	}
	if err != nil {
		return domain.Ship{}, fmt.Errorf("error on querying ship with id '%d': %w", mmsi, err)
	}

	return *domain.NewShip(mmsi, name, latitude, longitude), nil
}

func (pg *Postgres) Store(ctx context.Context, ships []domain.Ship) error {
	start := time.Now()
	defer pg.metrics.DBQueryTime("store_ship_data", start)

	for _, ship := range ships {
		_, err := pg.pool.Exec(ctx, updateSQL, ship.MMSI, ship.Name, ship.Latitude, ship.Longitude)
		if err != nil {
			return fmt.Errorf("error on inserting ship with mmsi '%d': %w", ship.MMSI, err)
		}
	}

	clog.Infof("Stored %d entries in Postgres in %d ms", len(ships), time.Since(start).Milliseconds())
	return nil
}

func (pg *Postgres) Shutdown(ctx context.Context) {
	pg.pool.Close()
}
