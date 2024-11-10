package repository_test

import (
	"context"
	"testing"

	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/pashagolub/pgxmock/v4"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/patraden/ya-practicum-go-shortly/internal/app/domain"
	e "github.com/patraden/ya-practicum-go-shortly/internal/app/errors"
	"github.com/patraden/ya-practicum-go-shortly/internal/app/logger"
	"github.com/patraden/ya-practicum-go-shortly/internal/app/repository"
)

func TestAddURLMapping(t *testing.T) {
	t.Parallel()

	log := logger.NewLogger(zerolog.InfoLevel).GetLogger()
	mockPool, err := pgxmock.NewPool()
	require.NoError(t, err)

	repo := repository.NewDBURLRepository(mockPool, log)
	ctx := context.Background()

	urlMapping := domain.NewURLMapping("a", "b")

	mockPool.
		ExpectExec(`INSERT INTO shortener.urlmapping \(slug, original, created_at, expires_at\)`).
		WithArgs(urlMapping.Slug, urlMapping.OriginalURL, urlMapping.CreatedAt, urlMapping.ExpiresAt).
		WillReturnResult(pgxmock.NewResult("INSERT", 1))

	err = repo.AddURLMapping(ctx, urlMapping)
	require.NoError(t, err)

	// collisions when duplicate
	mockPool.
		ExpectExec(`INSERT INTO shortener.urlmapping \(slug, original, created_at, expires_at\)`).
		WithArgs(urlMapping.Slug, urlMapping.OriginalURL, urlMapping.CreatedAt, urlMapping.ExpiresAt).
		WillReturnError(&pgconn.PgError{Code: pgerrcode.UniqueViolation})

	err = repo.AddURLMapping(ctx, urlMapping)
	require.ErrorIs(t, err, e.ErrRepoExists)
}

func TestRetriable(t *testing.T) {
	t.Parallel()

	log := logger.NewLogger(zerolog.InfoLevel).GetLogger()
	mockPool, err := pgxmock.NewPool()
	require.NoError(t, err)

	repo := repository.NewDBURLRepository(mockPool, log)
	ctx := context.Background()
	urlMapping := domain.NewURLMapping("a", "b")

	mockPool.
		ExpectExec(`INSERT INTO shortener.urlmapping \(slug, original, created_at, expires_at\)`).
		WithArgs(urlMapping.Slug, urlMapping.OriginalURL, urlMapping.CreatedAt, urlMapping.ExpiresAt).
		WillReturnError(&pgconn.PgError{Code: pgerrcode.ConnectionFailure})

	err = repo.AddURLMapping(ctx, urlMapping)
	require.ErrorContains(t, err, `call to method Exec() was not expected`)

	err = mockPool.ExpectationsWereMet()
	require.NoError(t, err)
}

func TestGetURLMapping(t *testing.T) {
	t.Parallel()

	log := logger.NewLogger(zerolog.InfoLevel).GetLogger()
	mockPool, err := pgxmock.NewPool()
	require.NoError(t, err)

	repo := repository.NewDBURLRepository(mockPool, log)
	ctx := context.Background()
	urlMapping := domain.NewURLMapping("a", "b")

	expctedRes := pgxmock.NewRows([]string{`slug`, `original`, `created_at`, `expires_at`}).
		AddRow(
			urlMapping.Slug,
			urlMapping.OriginalURL,
			urlMapping.CreatedAt,
			urlMapping.ExpiresAt,
		)

	mockPool.
		ExpectQuery(`SELECT slug, original, created_at, expires_at`).
		WithArgs(urlMapping.Slug).
		WillReturnRows(expctedRes)

	res, err := repo.GetURLMapping(ctx, urlMapping.Slug)
	require.NoError(t, err)

	assert.Equal(t, urlMapping.OriginalURL, res.OriginalURL)
}
