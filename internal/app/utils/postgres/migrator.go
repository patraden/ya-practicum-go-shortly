package postgres

import (
	"database/sql"

	// Import pgx driver for SQL compatibility.
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"

	e "github.com/patraden/ya-practicum-go-shortly/internal/app/errors"
)

type Migrator struct {
	db      *sql.DB
	dialect string
	dir     string
}

func NewMigrator(db *sql.DB, dialect, dir string) *Migrator {
	return &Migrator{
		db:      db,
		dialect: dialect,
		dir:     dir,
	}
}

func (m *Migrator) Up() error {
	if err := goose.SetDialect(m.dialect); err != nil {
		return e.Wrap("failed to set dialect:", err)
	}

	if err := goose.Up(m.db, m.dir); err != nil {
		return e.Wrap("failed to apply migrations:", err)
	}

	return nil
}

func NewPGMigrator(connString, dir string) (*Migrator, error) {
	db, err := sql.Open("pgx", connString)
	if err != nil {
		return nil, e.Wrap("failed to open db:", err)
	}

	return &Migrator{
		db:      db,
		dialect: "postgres",
		dir:     dir,
	}, nil
}
