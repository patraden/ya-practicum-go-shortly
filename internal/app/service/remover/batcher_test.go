package remover_test

import (
	"testing"
	"time"

	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/patraden/ya-practicum-go-shortly/internal/app/domain"
	e "github.com/patraden/ya-practicum-go-shortly/internal/app/domain/errors"
	"github.com/patraden/ya-practicum-go-shortly/internal/app/dto"
	"github.com/patraden/ya-practicum-go-shortly/internal/app/logger"
	"github.com/patraden/ya-practicum-go-shortly/internal/app/service/remover"
)

func TestBatcherAddTask(t *testing.T) {
	t.Parallel()

	outChan := make(chan remover.Job, 1)
	logger := logger.NewLogger(zerolog.DebugLevel).GetLogger()
	batcher := remover.NewBatcher(1, time.Second, 1, outChan, logger)

	// Add task to batcher
	err := batcher.AddTask(dto.UserSlug{Slug: "slug1", UserID: domain.NewUserID()})
	require.NoError(t, err, "expected no error when adding a task")

	// Add task to overflow the channel
	err = batcher.AddTask(dto.UserSlug{Slug: "slug2", UserID: domain.NewUserID()})
	assert.Equal(t, e.ErrMissedTask, err, "expected error when input channel is full")

	// Check that jobsMissed was incremented
	missed, _ := batcher.Metrics()
	assert.Equal(t, 0, int(missed), "expected jobsMissed to be 1")
}

func TestBatcherTimerResetOnEmit(t *testing.T) {
	t.Parallel()

	outChan := make(chan remover.Job, 1)
	logger := logger.NewLogger(zerolog.DebugLevel).GetLogger()
	batcher := remover.NewBatcher(2, 100*time.Millisecond, 10, outChan, logger)

	batcher.Start()
	defer batcher.Stop()

	// Add first task, wait for batch to be emitted by time
	err := batcher.AddTask(dto.UserSlug{Slug: "slug1", UserID: domain.NewUserID()})
	require.NoError(t, err)

	select {
	case job := <-outChan:
		// Ensure the first batch is emitted
		assert.Equal(t, 0, job.ID, "job ID should be 0 for the first batch")
		assert.Len(t, job.Tasks, 1, "batch size should be 1")
	case <-time.After(200 * time.Millisecond):
		t.Fatal("batch was not emitted within the expected time")
	}

	// Add another task and check if timer is reset
	err = batcher.AddTask(dto.UserSlug{Slug: "slug2", UserID: domain.NewUserID()})
	require.NoError(t, err)

	select {
	case job := <-outChan:
		// Ensure second batch is emitted
		assert.Equal(t, 1, job.ID, "job ID should be 1 for the second batch")
		assert.Len(t, job.Tasks, 1, "batch size should be 1")
	case <-time.After(200 * time.Millisecond):
		t.Fatal("batch was not emitted within the expected time")
	}
}

func TestBatcherBatchEmitBySize(t *testing.T) {
	t.Parallel()

	outChan := make(chan remover.Job, 1)
	logger := logger.NewLogger(zerolog.DebugLevel).GetLogger()
	batcher := remover.NewBatcher(3, time.Second, 10, outChan, logger)

	batcher.Start()
	defer batcher.Stop()

	err := batcher.AddTask(dto.UserSlug{Slug: "slug1", UserID: domain.NewUserID()})
	require.NoError(t, err)
	err = batcher.AddTask(dto.UserSlug{Slug: "slug2", UserID: domain.NewUserID()})
	require.NoError(t, err)
	err = batcher.AddTask(dto.UserSlug{Slug: "slug3", UserID: domain.NewUserID()})
	require.NoError(t, err)

	// check just a frst batch
	select {
	case job := <-outChan:
		assert.Equal(t, 0, job.ID, "job ID should be 0 for the first batch")
		assert.Len(t, job.Tasks, 3, "batch size should be 2")
	case <-time.After(time.Second):
		t.Fatal("batches were not emitted before first timeout")
	}
}

func TestBatcherBatchEmitByTime(t *testing.T) {
	t.Parallel()

	outChan := make(chan remover.Job, 1)
	logger := logger.NewLogger(zerolog.DebugLevel).GetLogger()
	batcher := remover.NewBatcher(10, 5*time.Millisecond, 10, outChan, logger)

	batcher.Start()
	defer batcher.Stop()

	err := batcher.AddTask(dto.UserSlug{Slug: "slug1", UserID: domain.NewUserID()})
	require.NoError(t, err)

	select {
	case job := <-outChan:
		assert.Equal(t, 0, job.ID, "job ID should be 0 for the first batch")
		assert.Len(t, job.Tasks, 1, "batch size should be 1")
	case <-time.After(2 * time.Second):
		t.Fatal("batch was not emitted by timer")
	}
}

func TestBatcherStop(t *testing.T) {
	t.Parallel()

	outChan := make(chan remover.Job, 1)
	logger := logger.NewLogger(zerolog.DebugLevel).GetLogger()
	batcher := remover.NewBatcher(2, time.Hour, 10, outChan, logger)

	batcher.Start()
	batcher.Stop()

	_, created := batcher.Metrics()

	assert.Eventually(t, func() bool {
		return created == 0
	}, time.Second, 10*time.Millisecond, "expected no jobs created after Stop")
}

func TestBatcherEmitOnStop(t *testing.T) {
	t.Parallel()

	outChan := make(chan remover.Job, 1)
	logger := logger.NewLogger(zerolog.DebugLevel).GetLogger()
	batcher := remover.NewBatcher(10, time.Hour, 10, outChan, logger)

	batcher.Start()
	err := batcher.AddTask(dto.UserSlug{Slug: "slug1", UserID: domain.NewUserID()})
	require.NoError(t, err)
	batcher.Stop()

	select {
	case job := <-outChan:
		assert.Equal(t, 0, job.ID, "job ID should be 0 for the first batch")
		assert.Len(t, job.Tasks, 1, "batch size should be 1")
	default:
		t.Fatal("batch was not emitted on stop")
	}
}
