package postgres_test

import (
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/require"

	"github.com/patraden/ya-practicum-go-shortly/internal/app/utils/postgres"
)

func TestMigrator(t *testing.T) {
	t.Parallel()

	database, mock, err := sqlmock.New()
	require.NoError(t, err)

	mock.ExpectBegin()
	mock.ExpectExec(".*").WillReturnError(err)

	tests := []struct {
		name    string
		dir     string
		dialect string
	}{
		{
			"test 1",
			"../../../../migrations",
			"postgres",
		},
		{
			"test 2",
			"migrations",
			"postgres",
		},
		{
			"test 3",
			"../../../../migrations",
			"qwerty",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			migrator := postgres.NewMigrator(database, test.dialect, test.dir)
			err = migrator.Up()
			require.Error(t, err)
		})
	}

	err = mock.ExpectationsWereMet()
	require.NoError(t, err)
}
