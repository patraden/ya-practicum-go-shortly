package remover

import (
	"context"
	"errors"
	"sync"
	"sync/atomic"

	"github.com/rs/zerolog"
)

type WorkerPool struct {
	maxWorkers    int
	jobs          chan Job
	jobResults    chan JobResult
	stoppedChan   chan struct{}
	jobsMissed    int32
	jobsSuccessed int32
	jobsFailed    int32
	wFunc         WorkerFunc
	log           *zerolog.Logger
	stopOnce      sync.Once
	runOnce       sync.Once

	wgWork sync.WaitGroup // WaitGroup for workers
	wgRes  sync.WaitGroup // WaitGroup for result collection
}

func NewWorkerPool(maxWorkers, jobsBufferSize int, wFunc WorkerFunc, log *zerolog.Logger) *WorkerPool {
	return &WorkerPool{
		maxWorkers:  maxWorkers,
		jobs:        make(chan Job, maxWorkers*jobsBufferSize),
		jobResults:  make(chan JobResult, maxWorkers*jobsBufferSize),
		stoppedChan: make(chan struct{}),
		wFunc:       wFunc,
		log:         log,
	}
}

func (pool *WorkerPool) JobsChannel() chan<- Job {
	return pool.jobs
}

// worker executes tasks from jobs chan and submits outputs to results chan.
func (pool *WorkerPool) worker(id int) {
	defer pool.wgWork.Done()

	pool.log.Info().
		Int("workerID", id).
		Msg("Worker added to pool")

	for job := range pool.jobs {
		result := pool.wFunc(job)
		select {
		case pool.jobResults <- result:
		default:
			atomic.AddInt32(&pool.jobsMissed, 1)
			pool.log.Error().
				Int("workerID", id).
				Int("JobID", job.ID).
				Int("jobResultsChSize", len(pool.jobResults)).
				Msg("Results channel full; dropping result")
		}
	}
}

func (pool *WorkerPool) Start() {
	pool.runOnce.Do(func() {
		pool.dispatch()

		pool.log.Info().Msg("Worker pool started")
	})
}

func (pool *WorkerPool) dispatch() {
	// start workers
	for i := range pool.maxWorkers {
		pool.wgWork.Add(1)
		go pool.worker(i)
	}

	// start collecting results
	pool.wgRes.Add(1)
	go pool.collectResults()

	// wait for results
	go func() {
		// once all results are collected we can close
		defer close(pool.stoppedChan)
		pool.wgRes.Wait()

		pool.log.Info().
			Int32("succeeded", pool.jobsSuccessed).
			Int32("failed", pool.jobsFailed).
			Int32("missed", pool.jobsMissed).
			Msg("Worker pool stopped gracefully.")
	}()
}

// collect worker results.
func (pool *WorkerPool) collectResults() {
	defer pool.wgRes.Done()

	for res := range pool.jobResults {
		if res.Err != nil {
			pool.log.Error().
				Err(res.Err).
				Int("JobID", res.ID).
				Msg("job failed")

			atomic.AddInt32(&pool.jobsFailed, 1)

			continue
		}

		pool.log.Info().
			Int("JobID", res.ID).
			Msg("job executed successfully")
		atomic.AddInt32(&pool.jobsSuccessed, 1)
	}
}

func (pool *WorkerPool) Stop(ctx context.Context) {
	pool.stopOnce.Do(func() {
		// close jobs channel and trigger workers to stop
		close(pool.jobs)

		pool.wgWork.Wait()
		// once workers stopped we can close results chan
		close(pool.jobResults)

		select {
		case <-pool.stoppedChan:
			return
		case <-ctx.Done():
			if errors.Is(ctx.Err(), context.Canceled) {
				pool.log.Error().
					Err(ctx.Err()).
					Msg("Worker pool stop canceled")
			} else {
				pool.log.Error().
					Err(ctx.Err()).
					Msg("Worker pool stop timed out")
			}
		}
	})
}

func (pool *WorkerPool) Metrics() (int32, int32, int32) {
	failed := atomic.LoadInt32(&pool.jobsFailed)
	missed := atomic.LoadInt32(&pool.jobsMissed)
	successed := atomic.LoadInt32(&pool.jobsSuccessed)

	return failed, missed, successed
}
