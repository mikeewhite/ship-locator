package postgres

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"

	"github.com/mikeewhite/ship-locator/backend/internal/core/domain"
	"github.com/mikeewhite/ship-locator/backend/pkg/apperrors"
	"github.com/mikeewhite/ship-locator/backend/pkg/clog"
	"github.com/mikeewhite/ship-locator/backend/pkg/config"
)

type Postgres struct {
	conn *pgx.Conn
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

// TODO - use a connection pool - https://github.com/jackc/pgx/wiki/Getting-started-with-pgx#using-a-connection-pool
func NewPostgres(ctx context.Context, cfg config.Config) (*Postgres, error) {
	url := fmt.Sprintf("postgres://%s:%s@%s/%s",
		cfg.PostgresUsername, cfg.PostgresPassword, cfg.PostgresAddress, cfg.PostgresDBName)
	conn, err := pgx.Connect(ctx, url)
	if err != nil {
		return nil, fmt.Errorf("failed to establish connection to postgres: %w", err)
	}

	// check the connection is healthy
	err = conn.Ping(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to ping postgres: %w", err)
	}

	return &Postgres{conn: conn}, nil
}

func (pg *Postgres) Get(ctx context.Context, mmsi int32) (domain.Ship, error) {
	var name string
	var latitude float64
	var longitude float64

	err := pg.conn.QueryRow(ctx, selectSQL, mmsi).Scan(&name, &latitude, &longitude)
	if err == pgx.ErrNoRows {
		return domain.Ship{}, apperrors.NewNoShipFoundErr(mmsi)
	}
	if err != nil {
		return domain.Ship{}, fmt.Errorf("error on querying ship with id '%d': %w", mmsi, err)
	}

	return *domain.NewShip(mmsi, name, latitude, longitude), nil
}

func (pg *Postgres) Store(ctx context.Context, ships []domain.Ship) error {
	for _, ship := range ships {
		_, err := pg.conn.Exec(ctx, updateSQL, ship.MMSI, ship.Name, ship.Latitude, ship.Longitude)
		if err != nil {
			return fmt.Errorf("error on inserting ship with mmsi '%d': %w", ship.MMSI, err)
		}
	}
	return nil
}

func (pg *Postgres) Shutdown(ctx context.Context) {
	if err := pg.conn.Close(ctx); err != nil {
		clog.Errorf("failed to close postgres connection: %s", err.Error())
	}
}
