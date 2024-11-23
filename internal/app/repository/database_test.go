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

	userID := domain.NewUserID()
	repo := repository.NewDBURLRepository(mockPool, log)
	ctx := context.Background()

	urlMapping := domain.NewURLMapping("a", "b", userID)

	mockPool.
		ExpectQuery(`INSERT INTO shortener.urlmapping \(slug, original, user_id, created_at, expires_at\)`).
		WithArgs(urlMapping.Slug, urlMapping.OriginalURL, urlMapping.UserID, urlMapping.CreatedAt, urlMapping.ExpiresAt).
		WillReturnRows(pgxmock.NewRows([]string{"slug", "original", "user_id", "created_at", "expires_at"}).
			AddRow(urlMapping.Slug, urlMapping.OriginalURL, urlMapping.UserID, urlMapping.CreatedAt, urlMapping.ExpiresAt))

	result, err := repo.AddURLMapping(ctx, urlMapping)
	require.NoError(t, err)
	require.Equal(t, urlMapping.Slug, result.Slug)
	require.Equal(t, urlMapping.OriginalURL, result.OriginalURL)
	require.Equal(t, urlMapping.UserID, result.UserID)
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

	userID := domain.NewUserID()
	repo := repository.NewDBURLRepository(mockPool, log)
	ctx := context.Background()

	urlMP := domain.NewURLMapping("slug1", "url1", userID)
	urlMPD := domain.NewURLMapping("slug2", "url1", userID)

	// unique vialation for duplicate slug
	mockPool.
		ExpectQuery(`INSERT INTO shortener.urlmapping \(slug, original, user_id, created_at, expires_at\)`).
		WithArgs(urlMP.Slug, urlMP.OriginalURL, urlMP.UserID, urlMP.CreatedAt, urlMP.ExpiresAt).
		WillReturnError(&pgconn.PgError{Code: pgerrcode.UniqueViolation})

	_, err = repo.AddURLMapping(ctx, urlMP)
	require.ErrorIs(t, err, e.ErrSlugExists)

	// duplicate url will not trigger error but rather return existing slug
	mockPool.
		ExpectQuery(`INSERT INTO shortener.urlmapping \(slug, original, user_id, created_at, expires_at\)`).
		WithArgs(urlMPD.Slug, urlMPD.OriginalURL, urlMPD.UserID, urlMPD.CreatedAt, urlMPD.ExpiresAt).
		WillReturnRows(pgxmock.NewRows([]string{"slug", "original", "user_id", "created_at", "expires_at"}).
			AddRow(urlMP.Slug, urlMP.OriginalURL, urlMP.UserID, urlMP.CreatedAt, urlMP.ExpiresAt))

	_, err = repo.AddURLMapping(ctx, urlMPD)
	require.ErrorIs(t, err, e.ErrOriginalExists)

	err = mockPool.ExpectationsWereMet()
	require.NoError(t, err)
}

func TestAddURLMappingRetriable(t *testing.T) {
	t.Parallel()

	log := logger.NewLogger(zerolog.InfoLevel).GetLogger()
	mockPool, err := pgxmock.NewPool()
	require.NoError(t, err)

	userID := domain.NewUserID()
	repo := repository.NewDBURLRepository(mockPool, log)
	ctx := context.Background()
	urlMP := domain.NewURLMapping("a", "b", userID)

	mockPool.
		ExpectQuery(`INSERT INTO shortener.urlmapping \(slug, original, user_id, created_at, expires_at\)`).
		WithArgs(urlMP.Slug, urlMP.OriginalURL, urlMP.UserID, urlMP.CreatedAt, urlMP.ExpiresAt).
		WillReturnError(&pgconn.PgError{Code: pgerrcode.ConnectionFailure}) // First retry
	mockPool.
		ExpectQuery(`INSERT INTO shortener.urlmapping \(slug, original, user_id, created_at, expires_at\)`).
		WithArgs(urlMP.Slug, urlMP.OriginalURL, urlMP.UserID, urlMP.CreatedAt, urlMP.ExpiresAt).
		WillReturnError(&pgconn.PgError{Code: pgerrcode.ConnectionFailure}) // Second retry
	mockPool.
		ExpectQuery(`INSERT INTO shortener.urlmapping \(slug, original, user_id, created_at, expires_at\)`).
		WithArgs(urlMP.Slug, urlMP.OriginalURL, urlMP.UserID, urlMP.CreatedAt, urlMP.ExpiresAt).
		WillReturnRows(pgxmock.NewRows([]string{"slug", "original", "user_id", "created_at", "expires_at"}).
			AddRow(urlMP.Slug, urlMP.OriginalURL, urlMP.UserID, urlMP.CreatedAt, urlMP.ExpiresAt)) // Success on third try

	result, err := repo.AddURLMapping(ctx, urlMP)
	require.NoError(t, err)
	require.Equal(t, urlMP.Slug, result.Slug)
	require.Equal(t, urlMP.OriginalURL, result.OriginalURL)
	require.Equal(t, urlMP.CreatedAt, result.CreatedAt)
	require.Equal(t, urlMP.ExpiresAt, result.ExpiresAt)

	err = mockPool.ExpectationsWereMet()
	require.NoError(t, err)
}

func TestGetURLMapping(t *testing.T) {
	t.Parallel()

	log := logger.NewLogger(zerolog.InfoLevel).GetLogger()
	mockPool, err := pgxmock.NewPool()
	require.NoError(t, err)

	userID := domain.NewUserID()
	repo := repository.NewDBURLRepository(mockPool, log)
	ctx := context.Background()
	urlMP := domain.NewURLMapping("a", "b", userID)

	// success
	expectedRes := pgxmock.NewRows([]string{"slug", "original", "user_id", "created_at", "expires_at"}).
		AddRow(urlMP.Slug, urlMP.OriginalURL, urlMP.UserID, urlMP.CreatedAt, urlMP.ExpiresAt)

	mockPool.
		ExpectQuery(`SELECT slug, original, user_id, created_at, expires_at`).
		WithArgs(urlMP.Slug).
		WillReturnRows(expectedRes)

	res, err := repo.GetURLMapping(ctx, urlMP.Slug)
	require.NoError(t, err)
	assert.Equal(t, urlMP.OriginalURL, res.OriginalURL)
	assert.Equal(t, urlMP.UserID, res.UserID)
	assert.Equal(t, urlMP.CreatedAt, res.CreatedAt)
	assert.Equal(t, urlMP.ExpiresAt, res.ExpiresAt)

	// failure not found
	mockPool.
		ExpectQuery(`SELECT slug, original, user_id, created_at, expires_at`).
		WithArgs(urlMP.Slug).
		WillReturnError(sql.ErrNoRows)

	res, err = repo.GetURLMapping(ctx, urlMP.Slug)
	require.ErrorIs(t, err, e.ErrSlugNotFound)
	assert.Nil(t, res)

	// failure retriable
	mockPool.
		ExpectQuery(`SELECT slug, original, user_id, created_at, expires_at`).
		WithArgs(urlMP.Slug).
		WillReturnError(&pgconn.PgError{Code: pgerrcode.ConnectionFailure})

	res, err = repo.GetURLMapping(ctx, urlMP.Slug)
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

	userID := domain.NewUserID()
	repo := repository.NewDBURLRepository(mockPool, log)
	ctx := context.Background()

	batch := &[]domain.URLMapping{
		*domain.NewURLMapping("a", "x", userID),
		*domain.NewURLMapping("b", "y", userID),
		*domain.NewURLMapping("c", "z", userID),
	}

	mockPool.ExpectBegin()
	mockPool.
		ExpectCopyFrom(
			[]string{"shortener", "urlmapping"},
			[]string{"slug", "original", "user_id", "created_at", "expires_at"}).
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

	userID := domain.NewUserID()
	repo := repository.NewDBURLRepository(mockPool, log)
	ctx := context.Background()

	batch := &[]domain.URLMapping{
		*domain.NewURLMapping("a", "x", userID),
		*domain.NewURLMapping("b", "y", userID),
		*domain.NewURLMapping("c", "z", userID),
	}

	mockPool.ExpectBegin()
	mockPool.
		ExpectCopyFrom(
			[]string{"shortener", "urlmapping"},
			[]string{"slug", "original", "user_id", "created_at", "expires_at"}).
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
		ExpectCopyFrom(
			[]string{"shortener", "urlmapping"},
			[]string{"slug", "original", "user_id", "created_at", "expires_at"}).
		WillReturnResult(3)
	mockPool.ExpectCommit().WillReturnError(e.ErrTestGeneral)

	err = repo.AddURLMappingBatch(ctx, batch)
	require.NoError(t, err) // commit errors just logged and ignored

	err = mockPool.ExpectationsWereMet()
	require.NoError(t, err)
}

func TestDBGetUserURLMappings(t *testing.T) {
	t.Parallel()

	log := logger.NewLogger(zerolog.InfoLevel).GetLogger()
	mockPool, err := pgxmock.NewPool()
	require.NoError(t, err)

	userID := domain.NewUserID()
	repo := repository.NewDBURLRepository(mockPool, log)
	ctx := context.Background()
	// Define test data
	urlMappings := []domain.URLMapping{
		*domain.NewURLMapping("a", "url1", userID),
		*domain.NewURLMapping("b", "url2", userID),
		*domain.NewURLMapping("c", "url3", userID),
	}
	// Mock a successful query
	rows := pgxmock.NewRows([]string{"slug", "original", "user_id", "created_at", "expires_at"})
	for _, mapping := range urlMappings {
		rows.AddRow(mapping.Slug, mapping.OriginalURL, mapping.UserID, mapping.CreatedAt, mapping.ExpiresAt)
	}

	mockPool.ExpectQuery(`SELECT slug, original, user_id, created_at, expires_at`).
		WithArgs(userID).
		WillReturnRows(rows)
	// Execute the method and verify results
	result, err := repo.GetUserURLMappings(ctx, userID)
	require.NoError(t, err)
	assert.Len(t, result, len(urlMappings))

	for i, res := range result {
		assert.Equal(t, urlMappings[i].Slug, res.Slug)
		assert.Equal(t, urlMappings[i].OriginalURL, res.OriginalURL)
		assert.Equal(t, urlMappings[i].UserID, res.UserID)
		assert.Equal(t, urlMappings[i].CreatedAt, res.CreatedAt)
		assert.Equal(t, urlMappings[i].ExpiresAt, res.ExpiresAt)
	}
	// Mock a "no rows" scenario
	mockPool.ExpectQuery(`SELECT slug, original, user_id, created_at, expires_at`).
		WithArgs(userID).
		WillReturnError(sql.ErrNoRows)

	result, err = repo.GetUserURLMappings(ctx, userID)
	require.ErrorIs(t, err, e.ErrUserNotFound)
	assert.Empty(t, result)
	// Mock a retriable database error
	mockPool.ExpectQuery(`SELECT slug, original, user_id, created_at, expires_at`).
		WithArgs(userID).
		WillReturnError(&pgconn.PgError{Code: pgerrcode.ConnectionFailure})

	result, err = repo.GetUserURLMappings(ctx, userID)
	require.Error(t, err)
	assert.Empty(t, result)
	// Verify expectations
	err = mockPool.ExpectationsWereMet()
	require.NoError(t, err)
}
