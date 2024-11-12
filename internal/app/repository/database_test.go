package repository_test

import (
	"context"
	"database/sql"
	"testing"

	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/pashagolub/pgxmock/v4"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/patraden/ya-practicum-go-shortly/internal/app/domain"
	e "github.com/patraden/ya-practicum-go-shortly/internal/app/domain/errors"
	"github.com/patraden/ya-practicum-go-shortly/internal/app/logger"
	"github.com/patraden/ya-practicum-go-shortly/internal/app/repository"
)

func TestAddURLMappingSuccess(t *testing.T) {
	t.Parallel()

	log := logger.NewLogger(zerolog.InfoLevel).GetLogger()
	mockPool, err := pgxmock.NewPool()
	require.NoError(t, err)

	repo := repository.NewDBURLRepository(mockPool, log)
	ctx := context.Background()

	urlMapping := domain.NewURLMapping("a", "b")

	mockPool.
		ExpectQuery(`INSERT INTO shortener.urlmapping \(slug, original, created_at, expires_at\)`).
		WithArgs(urlMapping.Slug, urlMapping.OriginalURL, urlMapping.CreatedAt, urlMapping.ExpiresAt).
		WillReturnRows(pgxmock.NewRows([]string{"slug", "original", "created_at", "expires_at"}).
			AddRow(urlMapping.Slug, urlMapping.OriginalURL, urlMapping.CreatedAt, urlMapping.ExpiresAt))

	result, err := repo.AddURLMapping(ctx, urlMapping)
	require.NoError(t, err)
	require.Equal(t, urlMapping.Slug, result.Slug)
	require.Equal(t, urlMapping.OriginalURL, result.OriginalURL)
	require.Equal(t, urlMapping.CreatedAt, result.CreatedAt)
	require.Equal(t, urlMapping.ExpiresAt, result.ExpiresAt)

	err = mockPool.ExpectationsWereMet()
	require.NoError(t, err)
}

func TestAddURLMappingErrors(t *testing.T) {
	t.Parallel()

	log := logger.NewLogger(zerolog.InfoLevel).GetLogger()
	mockPool, err := pgxmock.NewPool()
	require.NoError(t, err)

	repo := repository.NewDBURLRepository(mockPool, log)
	ctx := context.Background()

	urlMapping := domain.NewURLMapping("slug1", "url1")
	urlMappingDup := domain.NewURLMapping("slug2", "url1")

	// unique vialation for duplicate slug
	mockPool.
		ExpectQuery(`INSERT INTO shortener.urlmapping \(slug, original, created_at, expires_at\)`).
		WithArgs(urlMapping.Slug, urlMapping.OriginalURL, urlMapping.CreatedAt, urlMapping.ExpiresAt).
		WillReturnError(&pgconn.PgError{Code: pgerrcode.UniqueViolation})

	_, err = repo.AddURLMapping(ctx, urlMapping)
	require.ErrorIs(t, err, e.ErrSlugExists)

	// duplicate url will not trigger error but rather return existing slug
	mockPool.
		ExpectQuery(`INSERT INTO shortener.urlmapping \(slug, original, created_at, expires_at\)`).
		WithArgs(urlMappingDup.Slug, urlMappingDup.OriginalURL, urlMappingDup.CreatedAt, urlMappingDup.ExpiresAt).
		WillReturnRows(pgxmock.NewRows([]string{"slug", "original", "created_at", "expires_at"}).
			AddRow(urlMapping.Slug, urlMapping.OriginalURL, urlMapping.CreatedAt, urlMapping.ExpiresAt))

	_, err = repo.AddURLMapping(ctx, urlMappingDup)
	require.ErrorIs(t, err, e.ErrOriginalExists)

	err = mockPool.ExpectationsWereMet()
	require.NoError(t, err)
}

func TestAddURLMappingRetriable(t *testing.T) {
	t.Parallel()

	log := logger.NewLogger(zerolog.InfoLevel).GetLogger()
	mockPool, err := pgxmock.NewPool()
	require.NoError(t, err)

	repo := repository.NewDBURLRepository(mockPool, log)
	ctx := context.Background()
	urlMapping := domain.NewURLMapping("a", "b")

	mockPool.
		ExpectQuery(`INSERT INTO shortener.urlmapping \(slug, original, created_at, expires_at\)`).
		WithArgs(urlMapping.Slug, urlMapping.OriginalURL, urlMapping.CreatedAt, urlMapping.ExpiresAt).
		WillReturnError(&pgconn.PgError{Code: pgerrcode.ConnectionFailure}) // First retry
	mockPool.
		ExpectQuery(`INSERT INTO shortener.urlmapping \(slug, original, created_at, expires_at\)`).
		WithArgs(urlMapping.Slug, urlMapping.OriginalURL, urlMapping.CreatedAt, urlMapping.ExpiresAt).
		WillReturnError(&pgconn.PgError{Code: pgerrcode.ConnectionFailure}) // Second retry
	mockPool.
		ExpectQuery(`INSERT INTO shortener.urlmapping \(slug, original, created_at, expires_at\)`).
		WithArgs(urlMapping.Slug, urlMapping.OriginalURL, urlMapping.CreatedAt, urlMapping.ExpiresAt).
		WillReturnRows(pgxmock.NewRows([]string{"slug", "original", "created_at", "expires_at"}).
			AddRow(urlMapping.Slug, urlMapping.OriginalURL, urlMapping.CreatedAt, urlMapping.ExpiresAt)) // Success on third try

	result, err := repo.AddURLMapping(ctx, urlMapping)
	require.NoError(t, err)
	require.Equal(t, urlMapping.Slug, result.Slug)
	require.Equal(t, urlMapping.OriginalURL, result.OriginalURL)
	require.Equal(t, urlMapping.CreatedAt, result.CreatedAt)
	require.Equal(t, urlMapping.ExpiresAt, result.ExpiresAt)

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

	// success
	expectedRes := pgxmock.NewRows([]string{"slug", "original", "created_at", "expires_at"}).
		AddRow(urlMapping.Slug, urlMapping.OriginalURL, urlMapping.CreatedAt, urlMapping.ExpiresAt)

	mockPool.
		ExpectQuery(`SELECT slug, original, created_at, expires_at`).
		WithArgs(urlMapping.Slug).
		WillReturnRows(expectedRes)

	res, err := repo.GetURLMapping(ctx, urlMapping.Slug)
	require.NoError(t, err)
	assert.Equal(t, urlMapping.OriginalURL, res.OriginalURL)
	assert.Equal(t, urlMapping.CreatedAt, res.CreatedAt)
	assert.Equal(t, urlMapping.ExpiresAt, res.ExpiresAt)

	// failure not found
	mockPool.
		ExpectQuery(`SELECT slug, original, created_at, expires_at`).
		WithArgs(urlMapping.Slug).
		WillReturnError(sql.ErrNoRows)

	res, err = repo.GetURLMapping(ctx, urlMapping.Slug)
	require.ErrorIs(t, err, e.ErrSlugNotFound)
	assert.Nil(t, res)

	// failure retriable
	mockPool.
		ExpectQuery(`SELECT slug, original, created_at, expires_at`).
		WithArgs(urlMapping.Slug).
		WillReturnError(&pgconn.PgError{Code: pgerrcode.ConnectionFailure})

	res, err = repo.GetURLMapping(ctx, urlMapping.Slug)
	require.Error(t, err)
	assert.Nil(t, res)

	err = mockPool.ExpectationsWereMet()
	require.NoError(t, err)
}

func TestAddURLMappingBatchSuccess(t *testing.T) {
	t.Parallel()

	log := logger.NewLogger(zerolog.InfoLevel).GetLogger()
	mockPool, err := pgxmock.NewPool()
	require.NoError(t, err)

	repo := repository.NewDBURLRepository(mockPool, log)
	ctx := context.Background()

	batch := &[]domain.URLMapping{
		*domain.NewURLMapping("a", "x"),
		*domain.NewURLMapping("b", "y"),
		*domain.NewURLMapping("c", "z"),
	}

	mockPool.ExpectBegin()
	mockPool.
		ExpectCopyFrom([]string{"shortener", "urlmapping"}, []string{"slug", "original", "created_at", "expires_at"}).
		WillReturnResult(3)
	mockPool.ExpectCommit()

	err = repo.AddURLMappingBatch(ctx, batch)
	require.NoError(t, err)

	err = mockPool.ExpectationsWereMet()
	require.NoError(t, err)
}

func TestAddURLMappingBatchFailure(t *testing.T) {
	t.Parallel()

	log := logger.NewLogger(zerolog.InfoLevel).GetLogger()
	mockPool, err := pgxmock.NewPool()
	require.NoError(t, err)

	repo := repository.NewDBURLRepository(mockPool, log)
	ctx := context.Background()

	batch := &[]domain.URLMapping{
		*domain.NewURLMapping("a", "x"),
		*domain.NewURLMapping("b", "y"),
		*domain.NewURLMapping("c", "z"),
	}

	mockPool.ExpectBegin()
	mockPool.
		ExpectCopyFrom([]string{"shortener", "urlmapping"}, []string{"slug", "original", "created_at", "expires_at"}).
		WillReturnError(e.ErrTestGeneral)
	mockPool.ExpectRollback()
	mockPool.ExpectCommit() // commit is done in any case

	err = repo.AddURLMappingBatch(ctx, batch)
	require.Error(t, err)
	require.Contains(t, err.Error(), "error while running batch tx")

	err = mockPool.ExpectationsWereMet()
	require.NoError(t, err)

	mockPool.ExpectBegin()
	mockPool.
		ExpectCopyFrom([]string{"shortener", "urlmapping"}, []string{"slug", "original", "created_at", "expires_at"}).
		WillReturnResult(3)
	mockPool.ExpectCommit().WillReturnError(e.ErrTestGeneral)

	err = repo.AddURLMappingBatch(ctx, batch)
	require.NoError(t, err) // commit errors just logged and ignored

	err = mockPool.ExpectationsWereMet()
	require.NoError(t, err)
}
