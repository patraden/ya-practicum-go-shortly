package remover

import (
	"context"
	"sync/atomic"
	"time"

	"github.com/rs/zerolog"

	"github.com/patraden/ya-practicum-go-shortly/internal/app/domain"
	e "github.com/patraden/ya-practicum-go-shortly/internal/app/domain/errors"
	"github.com/patraden/ya-practicum-go-shortly/internal/app/dto"
	"github.com/patraden/ya-practicum-go-shortly/internal/app/middleware"
	"github.com/patraden/ya-practicum-go-shortly/internal/app/repository"
)

const (
	workers         = 2
	jobsBuffer      = 1000
	jobsSize        = 50
	tasksBufferSize = 2 * jobsBuffer
	jobInterval     = 5 * time.Second
)

type URLRemover interface {
	RemoveUserSlugs(ctx context.Context, slugs []domain.Slug) error
}

type AsyncRemover struct {
	repo    repository.URLRepository
	batcher *Batcher
	pool    *WorkerPool
	log     *zerolog.Logger

	running int32
}

func NewAsyncRemover(
	jobInterval time.Duration,
	repo repository.URLRepository,
	log *zerolog.Logger,
) *AsyncRemover {
	// it should be parent context but I have not time for that now
	ctx := context.Background()

	wfn := func(j Job) JobResult {
		err := repo.DelUserURLMappings(ctx, &j.Tasks)

		return JobResult{ID: j.ID, Err: err}
	}

	pool := NewWorkerPool(workers, jobsBuffer, wfn, log)
	jobCh := pool.JobsChannel()
	batcher := NewBatcher(jobsSize, jobInterval, tasksBufferSize, jobCh, log)

	return &AsyncRemover{
		repo:    repo,
		batcher: batcher,
		pool:    pool,
		log:     log,
		running: 0,
	}
}

func (r *AsyncRemover) IsRunning() bool {
	return r.running == 1
}

func (r *AsyncRemover) Start() {
	// batcher stops asyncroniously
	r.batcher.Start()
	// pool starts syncroniously
	r.pool.Start()
	atomic.SwapInt32(&r.running, 1)
}

func (r *AsyncRemover) Stop(ctx context.Context) {
	// batcher stops syncroniously
	r.batcher.Stop()
	// pool stops syncroniously
	r.pool.Stop(ctx)
	atomic.SwapInt32(&r.running, 0)
}

func (r *AsyncRemover) RemoveUserSlugs(ctx context.Context, slugs []domain.Slug) error {
	if !r.IsRunning() {
		r.log.Error().Msg("remover is not running")

		return e.ErrRemoverInternal
	}

	userID, ok := middleware.GetUserID(ctx)
	if !ok {
		r.log.Error().Msg("failed to get userID from context")

		return e.ErrRemoverInternal
	}

	for _, slug := range slugs {
		userSlug := dto.UserSlug{UserID: userID, Slug: slug}
		if err := r.batcher.AddTask(userSlug); err != nil {
			// fail fast for now
			r.log.Error().Err(err).Msg("failed to submit task to batcher")

			return err
		}
	}

	return nil
}
