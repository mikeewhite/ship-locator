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

// TODO create an in mem version of this repo to use in service

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

// TODO should id be mmsi?
func (pg *Postgres) Get(ctx context.Context, id int32) (domain.Ship, error) {
	var mmsi int32
	var name string
	var latitude float64
	var longitude float64

	// TODO move SQL to const
	err := pg.conn.QueryRow(ctx, "select mmsi, name, latitude, longitude from ships where mmsi=$1", id).Scan(&mmsi, &name, &latitude, &longitude)
	if err == pgx.ErrNoRows {
		return domain.Ship{}, apperrors.NewNoShipFoundErr(id)
	}
	if err != nil {
		return domain.Ship{}, fmt.Errorf("error on querying ship with id '%d': %w", id, err)
	}

	return *domain.NewShip(mmsi, name, latitude, longitude), nil
}

func (pg *Postgres) Store(ctx context.Context, ships []domain.Ship) error {
	// TODO - this needs to upsert

	// TODO - batch this up
	for _, ship := range ships {
		// TODO - this can be a const
		sqlStatement := `
			INSERT INTO ships (mmsi, name, latitude, longitude)
			VALUES ($1, $2, $3, $4)
			ON CONFLICT (mmsi)
			DO
				UPDATE SET name = EXCLUDED.name, latitude = EXCLUDED.latitude, longitude = EXCLUDED.longitude`
		_, err := pg.conn.Exec(ctx, sqlStatement, ship.MMSI, ship.Name, ship.Latitude, ship.Longitude)
		if err != nil {
			return fmt.Errorf("error on inserting ship with mmsi '%d': %w", ship.MMSI, err)
		}
	}

	// TODO - change the signature and return the ID of the inserted record?

	return nil
}

func (pg *Postgres) Shutdown(ctx context.Context) {
	if err := pg.conn.Close(ctx); err != nil {
		clog.Errorf("failed to close postgres connection: %s", err.Error())
	}
}
