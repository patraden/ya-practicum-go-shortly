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
	e "github.com/patraden/ya-practicum-go-shortly/internal/app/errors"
	q "github.com/patraden/ya-practicum-go-shortly/internal/app/repository/dbqueries"
	"github.com/patraden/ya-practicum-go-shortly/internal/app/utils"
	"github.com/patraden/ya-practicum-go-shortly/internal/app/utils/postgres"
)

const (
	queryRetryInterval  = 100 * time.Millisecond
	queryMaxElapsedTime = 5 * time.Second
)

type DBURLRepository struct {
	connPool postgres.ConnenctionPool
	queries  *q.Queries
	log      zerolog.Logger
}

func NewDBURLRepository(pool postgres.ConnenctionPool, log zerolog.Logger) *DBURLRepository {
	return &DBURLRepository{
		connPool: pool,
		queries:  q.New(pool),
		log:      log,
	}
}

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
					Msg("collision error")

				return backoff.Permanent(e.ErrRepoExists)
			}

			return backoff.Permanent(err)
		}

		return backoff.Permanent(err)
	}

	err := backoff.Retry(operation, backoff.WithContext(boff, ctx))
	if err != nil {
		return e.Wrap("retry error:", err)
	}

	return nil
}

func (repo *DBURLRepository) AddURLMapping(ctx context.Context, m *domain.URLMapping) error {
	retriableQuery := func() error {
		return repo.queries.AddURLMapping(ctx, q.AddURLMappingParams{
			Slug:      m.Slug,
			Original:  m.OriginalURL,
			CreatedAt: m.CreatedAt,
			ExpiresAt: m.ExpiresAt,
		})
	}

	err := repo.WithRetry(ctx, retriableQuery)
	if err != nil {
		return e.Wrap("failed to add urlmapping:", err)
	}

	return nil
}

func (repo *DBURLRepository) GetURLMapping(ctx context.Context, slug domain.Slug) (*domain.URLMapping, error) {
	var urlMap *domain.URLMapping

	retriableQuery := func() error {
		qm, err := repo.queries.GetURLMapping(ctx, slug)

		if errors.Is(err, sql.ErrNoRows) {
			return e.ErrRepoNotFound
		}

		if err != nil {
			return e.Wrap("query error", err)
		}

		urlMap = &domain.URLMapping{
			Slug:        qm.Slug,
			OriginalURL: qm.Original,
			CreatedAt:   qm.CreatedAt,
			ExpiresAt:   qm.ExpiresAt,
		}

		return nil
	}

	err := repo.WithRetry(ctx, retriableQuery)
	if err != nil {
		return nil, e.Wrap("failed to get urlmapping:", err)
	}

	return urlMap, nil
}

func (repo *DBURLRepository) AddURLMappingBatch(ctx context.Context, batch *[]domain.URLMapping) error {
	retriableQuery := func() error {
		trx, err := repo.connPool.BeginTx(ctx, pgx.TxOptions{})
		if err != nil {
			return e.Wrap("failed to start batch tx", err)
		}

		defer func() {
			if err := trx.Commit(ctx); err != nil {
				repo.log.
					Error().
					Err(err).
					Msg("failed to commit batch tx")
			}
		}()

		txQueries := repo.queries.WithTx(trx)
		batchParams := make([]q.AddURLMappingBatchCopyParams, len(*batch))

		for i, urlMapping := range *batch {
			batchParams[i] = q.AddURLMappingBatchCopyParams{
				Slug:      urlMapping.Slug,
				Original:  urlMapping.OriginalURL,
				CreatedAt: urlMapping.CreatedAt,
				ExpiresAt: urlMapping.ExpiresAt,
			}
		}

		rowsAffected, err := txQueries.AddURLMappingBatchCopy(ctx, batchParams)
		if err != nil {
			if err := trx.Rollback(ctx); err != nil {
				return e.Wrap("failed to commit batch tx", err)
			}

			return e.Wrap("error while running batch tx", err)
		}

		repo.log.
			Info().
			Int64("rows_affected", rowsAffected).
			Msg("loaded urlmappings in batch tx")

		return nil
	}

	err := repo.WithRetry(ctx, retriableQuery)
	if err != nil {
		return e.Wrap("failed to add urlmapping:", err)
	}

	return nil
}
