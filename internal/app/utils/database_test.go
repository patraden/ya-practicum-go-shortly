package utils_test

import (
	"database/sql"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/require"

	e "github.com/patraden/ya-practicum-go-shortly/internal/app/errors"
	"github.com/patraden/ya-practicum-go-shortly/internal/app/logger"
	"github.com/patraden/ya-practicum-go-shortly/internal/app/utils"
)

func TestNewDB(t *testing.T) {
	t.Parallel()

	log := logger.NewLogger(zerolog.InfoLevel).GetLogger()

	tests := []struct {
		name    string
		opener  utils.SQLDBOpener
		builder utils.DataSourceBuilder
		driver  string
		dsn     string
		wantErr error
	}{
		{
			"test 1",
			sql.Open,
			utils.PGDataSourceBuilder,
			"pgx",
			"postgres://postgres:postgres@localhost:5432/praktikum?sslmode=disable",
			nil,
		},
		{
			"test 2",
			sql.Open,
			utils.PGDataSourceBuilder,
			"badDriver",
			"postgres://postgres:postgres@localhost:5432/praktikum?sslmode=disable",
			e.ErrDBOpen,
		},
		{
			"test 3",
			sql.Open,
			utils.PGDataSourceBuilder,
			"pgx",
			"mysql:/bad_dsn",
			e.ErrDBDSNParse,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			_, err := utils.NewDatabase(test.opener, test.builder, log, test.driver, test.dsn)
			require.ErrorIs(t, err, test.wantErr)
		})
	}
}

func TestDBPing(t *testing.T) {
	t.Parallel()

	log := logger.NewLogger(zerolog.InfoLevel).GetLogger()
	mockedDB, mock, err := sqlmock.New(sqlmock.MonitorPingsOption(true))
	require.NoError(t, err)

	sqldb, err := utils.NewDatabase(
		func(_, _ string) (*sql.DB, error) { return mockedDB, nil },
		utils.PGDataSourceBuilder,
		log,
		`pgx`,
		"postgres://localhost:5432/praktikum",
	)

	tests := []struct {
		name    string
		pingErr error
		wantErr error
	}{
		{"test 1", nil, nil},
		{"test 2", e.ErrTest, e.ErrDBPing},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			mock.ExpectPing().WillReturnError(test.pingErr)

			require.NoError(t, err)

			err = sqldb.Ping()
			require.ErrorIs(t, err, test.wantErr)
		})
	}

	mock.ExpectClose()

	err = sqldb.Close()
	require.NoError(t, err)

	err = mock.ExpectationsWereMet()
	require.NoError(t, err)
}
