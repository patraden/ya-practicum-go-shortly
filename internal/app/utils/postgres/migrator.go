package postgres

import (
	"database/sql"

	// Import pgx driver for SQL compatibility.
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"

	e "github.com/patraden/ya-practicum-go-shortly/internal/app/domain/errors"
)

// Migrator represents a migration manager for applying migrations to a PostgreSQL database.
type Migrator struct {
	db      *sql.DB
	dialect string
	dir     string
}

// NewMigrator creates a new Migrator instance for managing migrations.
func NewMigrator(db *sql.DB, dialect, dir string) *Migrator {
	return &Migrator{
		db:      db,
		dialect: dialect,
		dir:     dir,
	}
}

// Up applies migrations from the migration directory to the database.
func (m *Migrator) Up() error {
	if err := goose.SetDialect(m.dialect); err != nil {
		return e.Wrap("failed to set dialect:", err, errLabel)
	}

	if err := goose.Up(m.db, m.dir); err != nil {
		return e.Wrap("failed to apply migrations:", err, errLabel)
	}

	return nil
}

// NewPGMigrator creates a new Migrator instance using a PostgreSQL connection string and migration directory.
func NewPGMigrator(connString, dir string) (*Migrator, error) {
	db, err := sql.Open("pgx", connString)
	if err != nil {
		return nil, e.Wrap("failed to open db:", err, errLabel)
	}

	return &Migrator{
		db:      db,
		dialect: "postgres",
		dir:     dir,
	}, nil
}
