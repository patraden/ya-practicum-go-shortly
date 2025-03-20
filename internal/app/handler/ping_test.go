package handler_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/pashagolub/pgxmock/v4"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/require"

	"github.com/patraden/ya-practicum-go-shortly/internal/app/config"
	"github.com/patraden/ya-practicum-go-shortly/internal/app/handler"
	"github.com/patraden/ya-practicum-go-shortly/internal/app/logger"
	"github.com/patraden/ya-practicum-go-shortly/internal/app/utils/postgres"
)

func TestHandleDBPing(t *testing.T) {
	t.Parallel()

	log := logger.NewLogger(zerolog.Disabled).GetLogger()
	mockPool, err := pgxmock.NewPool()
	require.NoError(t, err)

	defer mockPool.Close()

	mockPool.ExpectPing().WillReturnError(nil)

	cfg := &config.Config{DatabaseDSN: "postgres://postgres:postgres@localhost:5432/praktikum?sslmode=disable"}
	db := postgres.NewDatabase(log, cfg.DatabaseDSN).WithPool(mockPool)
	pingHandler := handler.NewPingHandler(db, cfg, log)

	req := httptest.NewRequest(http.MethodGet, "/ping", nil)
	rec := httptest.NewRecorder()

	pingHandler.HandleDBPing(rec, req)

	require.Equal(t, http.StatusOK, rec.Code)
	require.NoError(t, mockPool.ExpectationsWereMet())
}

func TestHandleDBPingFail(t *testing.T) {
	t.Parallel()

	log := logger.NewLogger(zerolog.Disabled).GetLogger()
	mockPool, err := pgxmock.NewPool()
	require.NoError(t, err)

	defer mockPool.Close()

	mockPool.ExpectPing().WillReturnError(context.DeadlineExceeded)

	cfg := &config.Config{DatabaseDSN: "postgres://postgres:postgres@localhost:5432/praktikum?sslmode=disable"}
	db := postgres.NewDatabase(log, cfg.DatabaseDSN).WithPool(mockPool)
	pingHandler := handler.NewPingHandler(db, cfg, log)

	req := httptest.NewRequest(http.MethodGet, "/ping", nil)
	rec := httptest.NewRecorder()

	pingHandler.HandleDBPing(rec, req)

	require.Equal(t, http.StatusInternalServerError, rec.Code)
	require.Contains(t, rec.Body.String(), "database is not reachable")
	require.NoError(t, mockPool.ExpectationsWereMet())
}
