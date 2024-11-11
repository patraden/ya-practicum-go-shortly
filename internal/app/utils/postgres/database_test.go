package postgres_test

import (
	"context"
	"testing"

	"github.com/pashagolub/pgxmock/v4"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/require"

	"github.com/patraden/ya-practicum-go-shortly/internal/app/logger"
	"github.com/patraden/ya-practicum-go-shortly/internal/app/utils/postgres"
)

func TestPostgresDB(t *testing.T) {
	t.Parallel()

	log := logger.NewLogger(zerolog.InfoLevel).GetLogger()
	mockPool, err := pgxmock.NewPool()
	require.NoError(t, err)
	mockPool.ExpectPing().Times(2)

	tests := []struct {
		name        string
		dsn         string
		mocked      bool
		wantInitErr bool
		wantPingErr bool
	}{
		{
			"test 1",
			"postgres://postgres:postgres@localhost:5432/praktikum?sslmode=disable",
			true,
			false,
			false,
		},
		{
			"test 2",
			"bad_dsn",
			false,
			true,
			true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			db := postgres.NewDatabase(log, test.dsn)
			defer db.Close()

			if test.mocked {
				db = db.WithPool(mockPool)
			}

			err := db.Init(context.Background())
			if test.wantInitErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}

			err = db.Ping(context.Background())
			if test.wantPingErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}
