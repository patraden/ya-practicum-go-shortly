package remover

import (
	"sync"
	"sync/atomic"
	"time"

	"github.com/rs/zerolog"

	e "github.com/patraden/ya-practicum-go-shortly/internal/app/domain/errors"
	"github.com/patraden/ya-practicum-go-shortly/internal/app/dto"
)

type Batcher struct {
	maxBatchSize int               // Maximum number of tasks in a batch
	maxBatchTime time.Duration     // Maximum time before emitting a batch
	inChan       chan dto.UserSlug // Inbound tasks channel
	outChan      chan<- Job        // Outbound batched jobs channel
	jobsCreated  int32             // jobs counter
	jobsMissed   int32             // missed jobs counter
	logger       *zerolog.Logger

	mu           sync.Mutex
	currentBatch []dto.UserSlug
	timer        *time.Timer
	wg           sync.WaitGroup
}

func NewBatcher(
	maxBatchSize int,
	maxBatchTime time.Duration,
	bufferSize int,
	outChan chan<- Job,
	logger *zerolog.Logger,
) *Batcher {
	return &Batcher{
		maxBatchSize: maxBatchSize,
		maxBatchTime: maxBatchTime,
		inChan:       make(chan dto.UserSlug, bufferSize),
		outChan:      outChan,
		logger:       logger,
		timer:        time.NewTimer(maxBatchTime),
	}
}

func (b *Batcher) AddTask(task dto.UserSlug) error {
	select {
	case b.inChan <- task:
		return nil
	default:
		b.logger.
			Error().
			Str("Slug", task.Slug.String()).
			Err(e.ErrMissedTask).
			Msg("Task batcher input channel is full")

		return e.ErrMissedTask
	}
}

func (b *Batcher) Start() {
	b.wg.Add(1)
	go b.run()

	b.logger.Info().Msg("Task batcher started")
}

func (b *Batcher) Stop() {
	close(b.inChan)
	b.wg.Wait()
	b.logger.Info().
		Int("batch_size", len(b.currentBatch)).
		Msg("Batcher stopped gracefully")
}

func (b *Batcher) run() {
	defer b.wg.Done()

	for {
		select {
		case task, ok := <-b.inChan:
			if !ok {
				// Channel closed, process remaining tasks and exit
				if len(b.currentBatch) > 0 {
					b.emitBatch("Stop")
					b.resetBatch()
				}

				return
			}

			b.addTask(task)

			if len(b.currentBatch) >= b.maxBatchSize {
				b.emitBatch("Size")
				b.resetBatch()
				b.resetTimer()
			}
		case <-b.timer.C:
			if len(b.currentBatch) > 0 {
				b.emitBatch("Time")
				b.resetBatch()
			}

			b.resetTimer()
		}
	}
}

func (b *Batcher) addTask(task dto.UserSlug) {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.currentBatch = append(b.currentBatch, task)
}

func (b *Batcher) emitBatch(trigger string) {
	if len(b.currentBatch) == 0 {
		return
	}

	job := Job{
		ID:    int(atomic.AddInt32(&b.jobsCreated, 1) - 1),
		Tasks: b.currentBatch,
	}

	select {
	case b.outChan <- job:
		b.logger.Info().
			Int("JobID", job.ID).
			Str("Trigger", trigger).
			Int("batch_size", len(job.Tasks)).
			Msg("Batch job emitted")
	default:
		b.logger.Error().
			Int("JobID", job.ID).
			Str("Trigger", trigger).
			Int("batch_size", len(job.Tasks)).
			Msg("Batch jobs channel is full; dropping job")
		atomic.AddInt32(&b.jobsMissed, 1)
	}
}

func (b *Batcher) resetTimer() {
	b.timer.Reset(b.maxBatchTime)
}

func (b *Batcher) resetBatch() {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.currentBatch = b.currentBatch[:0]
}

func (b *Batcher) Metrics() (int32, int32) {
	missed := atomic.LoadInt32(&b.jobsMissed)
	created := atomic.LoadInt32(&b.jobsCreated)

	return missed, created
}
