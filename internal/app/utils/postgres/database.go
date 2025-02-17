package postgres

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog"

	e "github.com/patraden/ya-practicum-go-shortly/internal/app/domain/errors"
)

// Database represents a connection to the PostgreSQL database,
// encapsulating the connection pool and logging functionality.
type Database struct {
	connString string
	ConnPool   ConnenctionPool
	Log        *zerolog.Logger
}

// NewDatabase creates a new instance of the Database struct with the provided logger and connection string.
func NewDatabase(log *zerolog.Logger, connString string) *Database {
	return &Database{
		connString: connString,
		ConnPool:   nil,
		Log:        log,
	}
}

// WithPool assigns a pre-existing connection pool to the Database instance.
func (pg *Database) WithPool(pool ConnenctionPool) *Database {
	if pg.ConnPool != nil {
		pg.Log.Info().Msg("database connection will be replaced")
	}

	pg.ConnPool = pool

	return pg
}

// Init initializes the database connection pool by parsing the connection string and establishing a pool.
func (pg *Database) Init(ctx context.Context) error {
	if pg.ConnPool != nil {
		return nil
	}

	config, err := pgxpool.ParseConfig(pg.connString)
	if err != nil {
		return e.Wrap("failed to parse connection string", err, errLabel)
	}

	pool, err := pgxpool.NewWithConfig(ctx, config)
	if err != nil {
		return e.Wrap("failed to configure connection pool", err, errLabel)
	}

	pg.ConnPool = pool

	pg.Log.Info().Msg("database connections pool initialized")

	return nil
}

// Ping checks the database connection by attempting to ping the connection pool.
func (pg *Database) Ping(ctx context.Context) error {
	if pg.ConnPool == nil {
		return e.ErrPGEmptyPool
	}

	if err := pg.ConnPool.Ping(ctx); err != nil {
		return e.Wrap("failed to ping database", err, errLabel)
	}

	return nil
}

// Close closes the connection pool and logs the disconnection.
func (pg *Database) Close() {
	if pg.ConnPool == nil {
		return
	}

	pg.ConnPool.Close()
	pg.ConnPool = nil

	pg.Log.Info().Msg("disconnected from database pool")
}
