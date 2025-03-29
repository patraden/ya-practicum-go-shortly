package repository

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/cenkalti/backoff/v4"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/rs/zerolog"

	"github.com/patraden/ya-practicum-go-shortly/internal/app/domain"
	e "github.com/patraden/ya-practicum-go-shortly/internal/app/domain/errors"
	"github.com/patraden/ya-practicum-go-shortly/internal/app/dto"
	"github.com/patraden/ya-practicum-go-shortly/internal/app/memento"
	q "github.com/patraden/ya-practicum-go-shortly/internal/app/repository/dbqueries"
	"github.com/patraden/ya-practicum-go-shortly/internal/app/utils"
	"github.com/patraden/ya-practicum-go-shortly/internal/app/utils/postgres"
)

const (
	queryRetryInterval  = 100 * time.Millisecond
	queryMaxElapsedTime = 5 * time.Second
)

// DBURLRepository is responsible for interacting with the database to handle URL mappings.
type DBURLRepository struct {
	connPool postgres.ConnenctionPool
	queries  *q.Queries
	log      *zerolog.Logger
}

// NewDBURLRepository creates a new instance of DBURLRepository with a connection pool and logger.
func NewDBURLRepository(pool postgres.ConnenctionPool, log *zerolog.Logger) *DBURLRepository {
	return &DBURLRepository{
		connPool: pool,
		queries:  q.New(pool),
		log:      log,
	}
}

// WithRetry retries the execution of the provided query function in case of transient errors
// such as connection or query execution issues.
func (repo *DBURLRepository) WithRetry(ctx context.Context, query func() error) error {
	boff := utils.LinearBackoff(queryMaxElapsedTime, queryRetryInterval)

	operation := func() error {
		err := query()
		if err == nil {
			return nil
		}

		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			switch pgErr.Code {
			case
				// retryable errors
				pgerrcode.ConnectionException,
				pgerrcode.ConnectionDoesNotExist,
				pgerrcode.ConnectionFailure,
				pgerrcode.CannotConnectNow,
				pgerrcode.SQLClientUnableToEstablishSQLConnection,
				pgerrcode.TransactionResolutionUnknown:
				repo.log.
					Info().
					Err(err).
					Msg("retrying query after retirable error")

				return err
			case
				// permanent errors
				pgerrcode.UniqueViolation:
				repo.log.
					Info().
					Err(err).
					Msg("slug collision error")

				return backoff.Permanent(e.ErrSlugExists)
			}

			return backoff.Permanent(err)
		}

		return backoff.Permanent(err)
	}

	err := backoff.Retry(operation, backoff.WithContext(boff, ctx))
	if err != nil {
		return e.Wrap("retry error:", err, errLabel)
	}

	return nil
}

// AddURLMapping adds a new URL mapping to the database and returns the created URL mapping.
func (repo *DBURLRepository) AddURLMapping(ctx context.Context, urlMap *domain.URLMapping) (*domain.URLMapping, error) {
	var res *domain.URLMapping

	retriableQuery := func() error {
		qmp, err := repo.queries.AddURLMapping(ctx, q.AddURLMappingParams{
			Slug:      urlMap.Slug,
			Original:  urlMap.OriginalURL,
			UserID:    urlMap.UserID,
			CreatedAt: urlMap.CreatedAt,
			ExpiresAt: urlMap.ExpiresAt,
			Deleted:   urlMap.Deleted,
		})
		if err != nil {
			return e.Wrap("failed to query", err, errLabel)
		}

		res = &domain.URLMapping{
			Slug:        qmp.Slug,
			OriginalURL: qmp.Original,
			UserID:      qmp.UserID,
			CreatedAt:   qmp.CreatedAt,
			ExpiresAt:   qmp.ExpiresAt,
			Deleted:     qmp.Deleted,
		}

		return nil
	}

	err := repo.WithRetry(ctx, retriableQuery)
	if err != nil {
		return res, e.Wrap("failed to add urlmapping", err, errLabel)
	}

	if res.Slug != urlMap.Slug {
		return res, e.ErrOriginalExists
	}

	return res, nil
}

// GetURLMapping retrieves a URL mapping by its slug from the database.
func (repo *DBURLRepository) GetURLMapping(ctx context.Context, slug domain.Slug) (*domain.URLMapping, error) {
	var urlMap *domain.URLMapping

	retriableQuery := func() error {
		qmr, err := repo.queries.GetURLMapping(ctx, slug)

		if errors.Is(err, sql.ErrNoRows) {
			return e.ErrSlugNotFound
		}

		if err != nil {
			return e.Wrap("failed to query", err, errLabel)
		}

		urlMap = &domain.URLMapping{
			Slug:        qmr.Slug,
			OriginalURL: qmr.Original,
			UserID:      qmr.UserID,
			CreatedAt:   qmr.CreatedAt,
			ExpiresAt:   qmr.ExpiresAt,
			Deleted:     qmr.Deleted,
		}

		return nil
	}

	err := repo.WithRetry(ctx, retriableQuery)
	if err != nil {
		return nil, e.Wrap("failed to get urlmapping", err, errLabel)
	}

	return urlMap, nil
}

// GetUserURLMappings retrieves all URL mappings for a given user from the database.
func (repo *DBURLRepository) GetUserURLMappings(ctx context.Context, user domain.UserID) ([]domain.URLMapping, error) {
	var results []domain.URLMapping

	retriableQuery := func() error {
		qresults, err := repo.queries.GetUserURLMappings(ctx, user)

		if errors.Is(err, sql.ErrNoRows) || len(qresults) == 0 {
			return e.ErrUserNotFound
		}

		if err != nil {
			return e.Wrap("failed to query", err, errLabel)
		}

		results = make([]domain.URLMapping, len(qresults))
		for i, qm := range qresults {
			results[i] = domain.URLMapping{
				Slug:        qm.Slug,
				OriginalURL: qm.Original,
				UserID:      qm.UserID,
				CreatedAt:   qm.CreatedAt,
				ExpiresAt:   qm.ExpiresAt,
				Deleted:     qm.Deleted,
			}
		}

		return nil
	}

	err := repo.WithRetry(ctx, retriableQuery)
	if err != nil {
		return []domain.URLMapping{}, e.Wrap("failed to get user urlmappings", err, errLabel)
	}

	return results, nil
}

// AddURLMappingBatch adds multiple URL mappings in a single batch to the database.
func (repo *DBURLRepository) AddURLMappingBatch(ctx context.Context, batch *[]domain.URLMapping) error {
	retriableQuery := func() error {
		trx, err := repo.connPool.BeginTx(ctx, pgx.TxOptions{})
		if err != nil {
			return e.Wrap("failed to start batch tx", err, errLabel)
		}

		defer func() {
			if errCom := trx.Commit(ctx); errCom != nil {
				repo.log.
					Error().
					Err(errCom).
					Msg("failed to commit batch tx")
			}
		}()

		txQueries := repo.queries.WithTx(trx)
		batchParams := make([]q.AddURLMappingBatchCopyParams, len(*batch))

		for i, urlMapping := range *batch {
			batchParams[i] = q.AddURLMappingBatchCopyParams{
				Slug:      urlMapping.Slug,
				Original:  urlMapping.OriginalURL,
				UserID:    urlMapping.UserID,
				CreatedAt: urlMapping.CreatedAt,
				ExpiresAt: urlMapping.ExpiresAt,
				Deleted:   urlMapping.Deleted,
			}
		}

		rowsAffected, err := txQueries.AddURLMappingBatchCopy(ctx, batchParams)
		if err != nil {
			if errRB := trx.Rollback(ctx); errRB != nil {
				repo.log.
					Error().
					Err(errRB).
					Msg("failed to rollback batch tx")

				return e.Wrap("failed to rollback batch tx", err, errLabel)
			}

			return e.Wrap("error while running batch tx", err, errLabel)
		}

		repo.log.
			Info().
			Int64("rows_affected", rowsAffected).
			Msg("loaded urlmappings in batch tx")

		return nil
	}

	err := repo.WithRetry(ctx, retriableQuery)
	if err != nil {
		return e.Wrap("failed to add urlmapping:", err, errLabel)
	}

	return nil
}

// DelUserURLMappings deletes URL mappings for a user based on their slugs.
func (repo *DBURLRepository) DelUserURLMappings(ctx context.Context, tasks []dto.UserSlug) error {
	var err error
	var rowsAffected int64

	retriableQuery := func() error {
		trx, beginErr := repo.connPool.BeginTx(ctx, pgx.TxOptions{})
		if beginErr != nil {
			return e.Wrap("failed to start batch tx", beginErr, errLabel)
		}

		defer func() {
			if err != nil {
				repo.log.Error().Err(err).
					Msg("rollback batch tx")

				if rollbackErr := trx.Rollback(ctx); rollbackErr != nil {
					repo.log.Error().Err(rollbackErr).
						Msg("failed to rollback batch tx")
				}
			}
		}()

		txQueries := repo.queries.WithTx(trx)
		deleteSlugParams := make([]q.FillDeletedSlugTempTableParams, len(tasks))

		for i, task := range tasks {
			deleteSlugParams[i] = q.FillDeletedSlugTempTableParams{
				Slug:   task.Slug,
				UserID: task.UserID,
			}
		}

		if err = txQueries.CreateDeletedSlugTempTable(ctx); err != nil {
			return e.Wrap("error creating temp table", err, errLabel)
		}

		if rowsAffected, err = txQueries.FillDeletedSlugTempTable(ctx, deleteSlugParams); err != nil {
			return e.Wrap("error filling temp table", err, errLabel)
		}

		if err = txQueries.DeleteSlugsInTarget(ctx); err != nil {
			return e.Wrap("error deleting slugs in target", err, errLabel)
		}

		if err = trx.Commit(ctx); err != nil {
			return e.Wrap("failed to commit batch tx", err, errLabel)
		}

		repo.log.Info().Int64("rows_affected", rowsAffected).
			Msg("url mappings deleted in batch tx")

		return nil
	}

	if err = repo.WithRetry(ctx, retriableQuery); err != nil {
		return e.Wrap("failed to delete user URL mappings", err, errLabel)
	}

	return nil
}

// GetStats retrieves repo statistics.
func (repo *DBURLRepository) GetStats(ctx context.Context) (*dto.RepoStats, error) {
	var stats *dto.RepoStats

	retriableQuery := func() error {
		qresults, err := repo.queries.GetStats(ctx)
		if err != nil {
			return err
		}

		stats = &dto.RepoStats{
			CountSlugs: qresults.Countslugs,
			CountUsers: qresults.Countusers,
		}

		return nil
	}

	if err := repo.WithRetry(ctx, retriableQuery); err != nil {
		return stats, err
	}

	return stats, nil
}

// CreateMemento creates a memento of the current state of the repository.
func (repo *DBURLRepository) CreateMemento() (*memento.Memento, error) {
	return nil, e.ErrStateNotmplemented
}

// RestoreMemento restores the state of the repository from the given memento.
func (repo *DBURLRepository) RestoreMemento(_ *memento.Memento) error {
	return e.ErrStateNotmplemented
}
