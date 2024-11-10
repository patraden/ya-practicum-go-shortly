package postgres

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog"

	e "github.com/patraden/ya-practicum-go-shortly/internal/app/errors"
)

type Database struct {
	connString string
	ConnPool   ConnenctionPool
	Log        zerolog.Logger
}

func NewDatabase(log zerolog.Logger, connString string) *Database {
	return &Database{
		connString: connString,
		ConnPool:   nil,
		Log:        log,
	}
}

func (pg *Database) WithPool(pool ConnenctionPool) *Database {
	if pg.ConnPool != nil {
		pg.Log.
			Info().
			Msg("database connection pool will be replaced")
	}

	pg.ConnPool = pool

	return pg
}

func (pg *Database) Init(ctx context.Context) error {
	if pg.ConnPool != nil {
		return nil
	}

	config, err := pgxpool.ParseConfig(pg.connString)
	if err != nil {
		return e.Wrap("failed to parse connection string", err)
	}

	pool, err := pgxpool.NewWithConfig(ctx, config)
	if err != nil {
		return e.Wrap("failed to configure connection pool", err)
	}

	pg.ConnPool = pool

	pg.Log.
		Info().
		Msg("database connections pool initialized")

	return nil
}

func (pg *Database) Ping(ctx context.Context) error {
	if pg.ConnPool == nil {
		return e.ErrDBEmptyPool
	}

	if err := pg.ConnPool.Ping(ctx); err != nil {
		return e.Wrap("failed to ping database", err)
	}

	return nil
}

func (pg *Database) Close() {
	if pg.ConnPool == nil {
		return
	}

	pg.ConnPool.Close()
	pg.ConnPool = nil

	pg.Log.
		Info().
		Msg("disconnected from database pool")
}
