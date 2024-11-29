package remover_test

import (
	"context"
	"strconv"
	"testing"
	"time"

	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"

	"github.com/patraden/ya-practicum-go-shortly/internal/app/domain"
	e "github.com/patraden/ya-practicum-go-shortly/internal/app/domain/errors"
	"github.com/patraden/ya-practicum-go-shortly/internal/app/dto"
	"github.com/patraden/ya-practicum-go-shortly/internal/app/logger"
	r "github.com/patraden/ya-practicum-go-shortly/internal/app/service/remover"
)

func TestWorkerPoolSuccess(t *testing.T) {
	t.Parallel()

	workerFunc := func(job r.Job) r.JobResult {
		return r.JobResult{ID: job.ID, Err: nil}
	}

	logger := logger.NewLogger(zerolog.DebugLevel).GetLogger()
	pool := r.NewWorkerPool(10, 10, workerFunc, logger)

	pool.Start()

	jobCount := 10
	for i := range jobCount {
		userID := domain.NewUserID()
		job := r.Job{
			ID: i,
			Tasks: []dto.UserSlug{{
				Slug:   domain.Slug("slug" + strconv.Itoa(i)),
				UserID: userID,
			}},
		}
		pool.JobsChannel() <- job
	}

	pool.Stop(context.Background())

	failed, missed, succeeded := pool.Metrics()
	actualProcessed := failed + missed + succeeded
	assert.Equal(t, jobCount, int(actualProcessed))
	assert.Equal(t, 0, int(missed))
}

func TestWorkerPoolFailure(t *testing.T) {
	t.Parallel()

	workerFunc := func(job r.Job) r.JobResult {
		// Simulating a failure for some jobs
		if job.ID%2 == 0 {
			return r.JobResult{ID: job.ID, Err: e.ErrTestGeneral}
		}

		return r.JobResult{ID: job.ID, Err: nil}
	}

	logger := logger.NewLogger(zerolog.DebugLevel).GetLogger()
	pool := r.NewWorkerPool(5, 5, workerFunc, logger)

	pool.Start()

	jobCount := 10
	for i := range jobCount {
		userID := domain.NewUserID()
		job := r.Job{
			ID: i,
			Tasks: []dto.UserSlug{{
				Slug:   domain.Slug("slug" + strconv.Itoa(i)),
				UserID: userID,
			}},
		}
		pool.JobsChannel() <- job
	}

	// Stop the pool and wait for completion
	pool.Stop(context.Background())

	// Validate metrics
	failed, missed, succeeded := pool.Metrics()
	assert.Equal(t, jobCount/2, int(succeeded)) // Half succeeded
	assert.Equal(t, jobCount/2, int(failed))    // Half failed
	assert.Equal(t, 0, int(missed))             // No missed
}

func TestWorkerPoolGracefulShutdown(t *testing.T) {
	t.Parallel()

	workerFunc := func(job r.Job) r.JobResult {
		time.Sleep(500 * time.Millisecond)

		return r.JobResult{ID: job.ID, Err: nil}
	}

	logger := logger.NewLogger(zerolog.DebugLevel).GetLogger()
	pool := r.NewWorkerPool(3, 3, workerFunc, logger)

	pool.Start()

	// Submit jobs with DeleteTask including UserID
	jobCount := 5
	for i := range jobCount {
		userID := domain.NewUserID() // Generating a new UserID
		job := r.Job{
			ID: i,
			Tasks: []dto.UserSlug{{
				Slug:   domain.Slug("slug" + strconv.Itoa(i)),
				UserID: userID,
			}},
		}
		pool.JobsChannel() <- job
	}

	// Simulate shutdown with timeout
	// 2 seconds should be enough to collect all results
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	pool.Stop(ctx)

	// Validate metrics
	failed, missed, succeeded := pool.Metrics()
	actualProcessed := failed + missed + succeeded

	assert.Equal(t, jobCount, int(actualProcessed)) // All jobs should be processed
	assert.Equal(t, 0, int(missed))                 // No missed jobs
}
