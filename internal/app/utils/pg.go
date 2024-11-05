package utils

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/jackc/pgx/v5/pgconn"
	// Import pgx driver for SQL compatibility.
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/rs/zerolog"
)

func NewDB(log zerolog.Logger, dsn string) (*sql.DB, error) {
	cfg, err := pgconn.ParseConfig(dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to parse dsn: %w", err)
	}

	ps := fmt.Sprintf("host=%s user=%s password=%s dbname=%s sslmode=disable",
		cfg.Host, cfg.User, cfg.Password, cfg.Database)

	inst, err := sql.Open("pgx", ps)
	if err != nil {
		return nil, fmt.Errorf("failed to init db: %w", err)
	}

	if err := inst.PingContext(context.Background()); err != nil {
		return nil, fmt.Errorf("failed to ping db: %w", err)
	}

	log.Info().
		Str("host", cfg.Host).
		Str("db", cfg.Database).
		Str("user", cfg.User).
		Msg("Connected to db")

	return inst, nil
}
