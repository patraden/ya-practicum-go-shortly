package utils

import (
	"database/sql"
	"fmt"

	"github.com/jackc/pgx/v5/pgconn"
	// Import pgx driver for SQL compatibility.
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/rs/zerolog"

	e "github.com/patraden/ya-practicum-go-shortly/internal/app/errors"
)

type DBConfig struct {
	Host string
	User string
	DB   string
}

type (
	SQLDBOpener       func(driverName, dataSourceName string) (*sql.DB, error)
	DataSourceBuilder func(dsn string) (string, *DBConfig, error)
)

type SQLDatabase struct {
	db       *sql.DB
	dbConfig *DBConfig
	log      zerolog.Logger
}

func NewDatabase(
	dbOpener SQLDBOpener,
	dataSourceBuilder DataSourceBuilder,
	log zerolog.Logger,
	driverName string,
	dsn string,
) (*SQLDatabase, error) {
	dataSourceName, dbConfig, err := dataSourceBuilder(dsn)
	if err != nil {
		return nil, e.ErrDBDSNParse
	}

	sqldb, err := dbOpener(driverName, dataSourceName)
	if err != nil {
		return nil, e.ErrDBOpen
	}

	return &SQLDatabase{
		db:       sqldb,
		dbConfig: dbConfig,
		log:      log,
	}, nil
}

func (sqldb *SQLDatabase) Ping() error {
	if err := sqldb.db.Ping(); err != nil {
		return e.ErrDBPing
	}

	sqldb.log.Info().
		Str("host", sqldb.dbConfig.Host).
		Str("db", sqldb.dbConfig.DB).
		Str("user", sqldb.dbConfig.User).
		Msg("db is alive!")

	return nil
}

func (sqldb *SQLDatabase) Close() error {
	if err := sqldb.db.Close(); err != nil {
		return e.ErrDBClose
	}

	return nil
}

func PGDataSourceBuilder(dsn string) (string, *DBConfig, error) {
	cfg, err := pgconn.ParseConfig(dsn)
	if err != nil {
		return ``, nil, e.ErrDBDSNParse
	}

	dataSourceName := fmt.Sprintf("host=%s user=%s password=%s dbname=%s sslmode=disable",
		cfg.Host, cfg.User, cfg.Password, cfg.Database)

	return dataSourceName, &DBConfig{
		Host: cfg.Host,
		User: cfg.User,
		DB:   cfg.Database,
	}, nil
}

func NewPG(dsn string, log zerolog.Logger) (*SQLDatabase, error) {
	sqldb, err := NewDatabase(sql.Open, PGDataSourceBuilder, log, "pgx", dsn)
	if err != nil {
		return nil, err
	}

	return sqldb, nil
}
