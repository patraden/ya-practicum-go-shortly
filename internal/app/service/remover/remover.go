package remover

import (
	"context"
	"sync"
	"time"

	"github.com/rs/zerolog"

	"github.com/patraden/ya-practicum-go-shortly/internal/app/domain"
	e "github.com/patraden/ya-practicum-go-shortly/internal/app/domain/errors"
	"github.com/patraden/ya-practicum-go-shortly/internal/app/dto"
	"github.com/patraden/ya-practicum-go-shortly/internal/app/middleware"
	"github.com/patraden/ya-practicum-go-shortly/internal/app/repository"
	b "github.com/patraden/ya-practicum-go-shortly/pkg/batcher"
)

const (
	batchBuffer  = 1000
	batchMaxSize = 100
	batchTimeout = time.Second
)

// URLRemover is an interface for removing user slugs from the repository.
type URLRemover interface {
	RemoveUserSlugs(ctx context.Context, slugs []domain.Slug) error
}

// BatchRemover is a concrete implementation of the URLRemover interface
// that handles the removal of user slugs in batches.
type BatchRemover struct {
	repo    repository.URLRepository
	batcher *b.Batcher
	log     *zerolog.Logger
	wg      *sync.WaitGroup
}

// NewBatchRemover creates a new instance of BatchRemover with the specified repository and logger.
func NewBatchRemover(repo repository.URLRepository, log *zerolog.Logger) (*BatchRemover, error) {
	commitFn := func(ctx context.Context, batch b.Batch) {
		slugs := make([]dto.UserSlug, 0, len(batch))

		for _, op := range batch {
			if slug, ok := op.Value.(dto.UserSlug); ok {
				slugs = append(slugs, slug)
			} else {
				op.SetError(e.ErrFailedCast)
			}
		}

		if len(slugs) == 0 {
			return
		}

		err := repo.DelUserURLMappings(ctx, slugs)
		batch.SetError(err)

		if err != nil {
			log.Error().Err(err).
				Int("size", len(batch)).
				Msg("remover: batch failed")
		}

		select {
		case <-ctx.Done():
			log.Info().
				Int("size", len(batch)).
				Msg("remover: batch cancelled")

			return
		default:
			break
		}
	}

	batcher, err := b.New(
		commitFn,
		b.WithBufferSize(batchBuffer),
		b.WithTimeout(batchTimeout),
		b.WithMaxSize(batchMaxSize),
		b.WithLogger(log),
	)
	if err != nil {
		return nil, e.ErrRemoverInitBatcher
	}

	return &BatchRemover{
		repo:    repo,
		batcher: batcher,
		log:     log,
		wg:      &sync.WaitGroup{},
	}, nil
}

// Start initiates the batch processing in a separate goroutine.
func (r *BatchRemover) Start(ctx context.Context) {
	r.wg.Add(1)

	go func() {
		defer r.wg.Done()
		r.batcher.Batch(ctx)
	}()
}

// Stop gracefully stops the batch processor, waiting for all tasks to complete.
func (r *BatchRemover) Stop(ctx context.Context) {
	done := make(chan struct{})
	go func() {
		r.wg.Wait()
		close(done)
	}()

	select {
	case <-done:
		r.log.Info().
			Msg("remover: stopped gracefully")
	case <-ctx.Done():
		r.log.Error().
			Msg("remover: shutdown timed out")
	}
}

// RemoveUserSlugs removes a list of user slugs asynchronously by batching the requests.
func (r *BatchRemover) RemoveUserSlugs(ctx context.Context, slugs []domain.Slug) error {
	userID, ok := middleware.GetUserID(ctx)
	if !ok {
		return e.ErrRemoverInternal
	}

	errCount := 0

	for _, slug := range slugs {
		userSlug := dto.UserSlug{UserID: userID, Slug: slug}
		if _, err := r.batcher.Send(ctx, userSlug); err != nil {
			errCount++
		}
	}

	if errCount > 0 {
		r.log.Error().
			Int("count", errCount).
			Msg("remover: some slugs missed from batch")

		return e.ErrRemoverInternal
	}

	// batch operation allows to block code here by calling op.Wait()
	// but since requirements specify async deletion, this is skipped
	return nil
}
