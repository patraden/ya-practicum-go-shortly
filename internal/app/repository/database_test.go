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
	"github.com/patraden/ya-practicum-go-shortly/internal/app/dto"
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

	urlm := domain.NewURLMapping("a", "b", userID)

	mockPool.
		ExpectQuery(`INSERT INTO shortener.urlmapping \(slug, original, user_id, created_at, expires_at, deleted\)`).
		WithArgs(urlm.Slug, urlm.OriginalURL, urlm.UserID, urlm.CreatedAt, urlm.ExpiresAt, urlm.Deleted).
		WillReturnRows(pgxmock.NewRows([]string{"slug", "original", "user_id", "created_at", "expires_at", "deleted"}).
			AddRow(urlm.Slug, urlm.OriginalURL, urlm.UserID, urlm.CreatedAt, urlm.ExpiresAt, urlm.Deleted))

	result, err := repo.AddURLMapping(ctx, urlm)
	require.NoError(t, err)
	require.Equal(t, urlm.Slug, result.Slug)
	require.Equal(t, urlm.OriginalURL, result.OriginalURL)
	require.Equal(t, urlm.UserID, result.UserID)
	require.Equal(t, urlm.CreatedAt, result.CreatedAt)
	require.Equal(t, urlm.ExpiresAt, result.ExpiresAt)
	require.Equal(t, urlm.Deleted, result.Deleted)

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

	urlm := domain.NewURLMapping("slug1", "url1", userID)
	urlmd := domain.NewURLMapping("slug2", "url1", userID)

	// unique vialation for duplicate slug
	mockPool.
		ExpectQuery(`INSERT INTO shortener.urlmapping \(slug, original, user_id, created_at, expires_at, deleted\)`).
		WithArgs(urlm.Slug, urlm.OriginalURL, urlm.UserID, urlm.CreatedAt, urlm.ExpiresAt, urlm.Deleted).
		WillReturnError(&pgconn.PgError{Code: pgerrcode.UniqueViolation})

	_, err = repo.AddURLMapping(ctx, urlm)
	require.ErrorIs(t, err, e.ErrSlugExists)

	// duplicate url will not trigger error but rather return existing slug
	mockPool.
		ExpectQuery(`INSERT INTO shortener.urlmapping \(slug, original, user_id, created_at, expires_at, deleted\)`).
		WithArgs(urlmd.Slug, urlmd.OriginalURL, urlmd.UserID, urlmd.CreatedAt, urlmd.ExpiresAt, urlmd.Deleted).
		WillReturnRows(pgxmock.NewRows([]string{"slug", "original", "user_id", "created_at", "expires_at", "deleted"}).
			AddRow(urlm.Slug, urlm.OriginalURL, urlm.UserID, urlm.CreatedAt, urlm.ExpiresAt, urlm.Deleted))

	_, err = repo.AddURLMapping(ctx, urlmd)
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
	urlm := domain.NewURLMapping("a", "b", userID)

	mockPool.
		ExpectQuery(`INSERT INTO shortener.urlmapping \(slug, original, user_id, created_at, expires_at, deleted\)`).
		WithArgs(urlm.Slug, urlm.OriginalURL, urlm.UserID, urlm.CreatedAt, urlm.ExpiresAt, urlm.Deleted).
		WillReturnError(&pgconn.PgError{Code: pgerrcode.ConnectionFailure}) // First retry
	mockPool.
		ExpectQuery(`INSERT INTO shortener.urlmapping \(slug, original, user_id, created_at, expires_at, deleted\)`).
		WithArgs(urlm.Slug, urlm.OriginalURL, urlm.UserID, urlm.CreatedAt, urlm.ExpiresAt, urlm.Deleted).
		WillReturnError(&pgconn.PgError{Code: pgerrcode.ConnectionFailure}) // Second retry
	// Success on third try
	mockPool.
		ExpectQuery(`INSERT INTO shortener.urlmapping \(slug, original, user_id, created_at, expires_at, deleted\)`).
		WithArgs(urlm.Slug, urlm.OriginalURL, urlm.UserID, urlm.CreatedAt, urlm.ExpiresAt, urlm.Deleted).
		WillReturnRows(pgxmock.NewRows([]string{"slug", "original", "user_id", "created_at", "expires_at", "deleted"}).
			AddRow(urlm.Slug, urlm.OriginalURL, urlm.UserID, urlm.CreatedAt, urlm.ExpiresAt, urlm.Deleted))

	result, err := repo.AddURLMapping(ctx, urlm)
	require.NoError(t, err)
	require.Equal(t, urlm.Slug, result.Slug)
	require.Equal(t, urlm.OriginalURL, result.OriginalURL)
	require.Equal(t, urlm.CreatedAt, result.CreatedAt)
	require.Equal(t, urlm.ExpiresAt, result.ExpiresAt)
	require.Equal(t, urlm.Deleted, result.Deleted)

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
	urlm := domain.NewURLMapping("a", "b", userID)

	// success
	expectedRes := pgxmock.NewRows([]string{"slug", "original", "user_id", "created_at", "expires_at", "deleted"}).
		AddRow(urlm.Slug, urlm.OriginalURL, urlm.UserID, urlm.CreatedAt, urlm.ExpiresAt, urlm.Deleted)

	mockPool.
		ExpectQuery(`SELECT slug, original, user_id, created_at, expires_at, deleted`).
		WithArgs(urlm.Slug).
		WillReturnRows(expectedRes)

	res, err := repo.GetURLMapping(ctx, urlm.Slug)
	require.NoError(t, err)
	assert.Equal(t, urlm.OriginalURL, res.OriginalURL)
	assert.Equal(t, urlm.UserID, res.UserID)
	assert.Equal(t, urlm.CreatedAt, res.CreatedAt)
	assert.Equal(t, urlm.ExpiresAt, res.ExpiresAt)
	assert.Equal(t, urlm.Deleted, res.Deleted)

	// failure not found
	mockPool.
		ExpectQuery(`SELECT slug, original, user_id, created_at, expires_at`).
		WithArgs(urlm.Slug).
		WillReturnError(sql.ErrNoRows)

	res, err = repo.GetURLMapping(ctx, urlm.Slug)
	require.ErrorIs(t, err, e.ErrSlugNotFound)
	assert.Nil(t, res)

	// failure retriable
	mockPool.
		ExpectQuery(`SELECT slug, original, user_id, created_at, expires_at`).
		WithArgs(urlm.Slug).
		WillReturnError(&pgconn.PgError{Code: pgerrcode.ConnectionFailure})

	res, err = repo.GetURLMapping(ctx, urlm.Slug)
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
			[]string{"slug", "original", "user_id", "created_at", "expires_at", "deleted"}).
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
			[]string{"slug", "original", "user_id", "created_at", "expires_at", "deleted"}).
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
			[]string{"slug", "original", "user_id", "created_at", "expires_at", "deleted"}).
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
	rows := pgxmock.NewRows([]string{"slug", "original", "user_id", "created_at", "expires_at", "deleted"})
	for _, m := range urlMappings {
		rows.AddRow(m.Slug, m.OriginalURL, m.UserID, m.CreatedAt, m.ExpiresAt, m.Deleted)
	}

	mockPool.ExpectQuery(`SELECT slug, original, user_id, created_at, expires_at, deleted`).
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
		assert.Equal(t, urlMappings[i].Deleted, res.Deleted)
	}
	// Mock a "no rows" scenario
	mockPool.ExpectQuery(`SELECT slug, original, user_id, created_at, expires_at, deleted`).
		WithArgs(userID).
		WillReturnError(sql.ErrNoRows)

	result, err = repo.GetUserURLMappings(ctx, userID)
	require.ErrorIs(t, err, e.ErrUserNotFound)
	assert.Empty(t, result)
	// Mock a retriable database error
	mockPool.ExpectQuery(`SELECT slug, original, user_id, created_at, expires_at, deleted`).
		WithArgs(userID).
		WillReturnError(&pgconn.PgError{Code: pgerrcode.ConnectionFailure})

	result, err = repo.GetUserURLMappings(ctx, userID)
	require.Error(t, err)
	assert.Empty(t, result)
	// Verify expectations
	err = mockPool.ExpectationsWereMet()
	require.NoError(t, err)
}

func TestDelUserURLMappingsSuccess(t *testing.T) {
	t.Parallel()

	log := logger.NewLogger(zerolog.InfoLevel).GetLogger()
	mockPool, err := pgxmock.NewPool()
	require.NoError(t, err)

	repo := repository.NewDBURLRepository(mockPool, log)
	ctx := context.Background()

	userSlugTasks := []dto.UserSlug{
		{Slug: "slug1", UserID: domain.NewUserID()},
		{Slug: "slug2", UserID: domain.NewUserID()},
	}

	mockPool.ExpectBegin()
	mockPool.ExpectExec(`CREATE TEMP TABLE urlmapping_tmp`).WillReturnResult(pgxmock.NewResult("CREATE", 0))
	mockPool.ExpectCopyFrom(
		[]string{"urlmapping_tmp"},
		[]string{"slug", "user_id"}).
		WillReturnResult(2)
	mockPool.ExpectExec(`UPDATE shortener.urlmapping`).WillReturnResult(pgxmock.NewResult("UPDATE", 0))
	mockPool.ExpectCommit()

	err = repo.DelUserURLMappings(ctx, &userSlugTasks)
	require.NoError(t, err)

	err = mockPool.ExpectationsWereMet()
	require.NoError(t, err)
}

func TestDelUserURLMappingsFailure(t *testing.T) {
	t.Parallel()

	log := logger.NewLogger(zerolog.InfoLevel).GetLogger()
	mockPool, err := pgxmock.NewPool()
	require.NoError(t, err)

	repo := repository.NewDBURLRepository(mockPool, log)
	ctx := context.Background()

	userSlugTasks := []dto.UserSlug{
		{Slug: "slug1", UserID: domain.NewUserID()},
	}

	mockPool.ExpectBegin()
	mockPool.ExpectExec(`CREATE TEMP TABLE urlmapping_tmp`).WillReturnResult(pgxmock.NewResult("CREATE", 0))
	mockPool.ExpectCopyFrom(
		[]string{"urlmapping_tmp"},
		[]string{"slug", "user_id"}).
		WillReturnError(e.ErrTestGeneral)
	mockPool.ExpectRollback()

	// Call the method under test
	err = repo.DelUserURLMappings(ctx, &userSlugTasks)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "error filling temp table")

	err = mockPool.ExpectationsWereMet()
	require.NoError(t, err)
}
