package batcher

import (
	"context"
	"os"
	"sync/atomic"
	"time"

	"github.com/rs/zerolog"
)

// Batcher aux constats.
const (
	UnlimitedSize               = 0    // Batch an unlimited amount of operations.
	NoTimeout     time.Duration = 0    // Batch operations for an infinite duration.
	DefaultBuffer               = 1000 // Default operations buffer size.
)

// Batcher configuration option.
type Option func(*Batcher)

// WithMaxSize configures the max size constraint on a batcher.
func WithMaxSize(maxSize int) Option {
	return func(b *Batcher) {
		b.maxSize = maxSize
	}
}

// WithTimeout configures the timeout constraint on a batcher.
func WithTimeout(timeout time.Duration) Option {
	return func(b *Batcher) {
		b.timeout = timeout
	}
}

// WithBufferSize configure buffer size for operations to improve performance.
func WithBufferSize(bufferSize int) Option {
	size := DefaultBuffer
	if bufferSize > 0 {
		size = bufferSize
	}

	return func(b *Batcher) {
		b.in = make(chan *Operation, size)
	}
}

// WithLogger configures zerolog.Logger to be used by Batcher.
func WithLogger(log *zerolog.Logger) Option {
	return func(b *Batcher) {
		b.log = log
	}
}

// Batcher provides a generic and versatile implementation of a batching algorithm for Golang,
// with no third-party dependencies.
//
// The algorithm can be constrained in space and time, with a simple yet robust API,
// enabling developers to easily incorporate batching into their live services.
type Batcher struct {
	commitFn CommitFunc
	maxSize  int
	timeout  time.Duration
	in       chan *Operation
	log      *zerolog.Logger
	closing  uint32
}

// New creates a new batcher, calling the commit function each time it
// completes a batch of operations according to its options. It panics if the
// commit function is nil, max size is negative, timeout is negative or max
// size equals [UnlimitedSize] and timeout equals [NoTimeout] (the default if
// no options are provided).
//
// Some examples:
//
// Create a batcher committing a batch every 10 operations:
//
//	New(commitFn, WithMaxSize(10))
//
// Create a batcher committing a batch every second:
//
//	New(commitFn, WithTimeout(1 * time.Second))
//
// Create a batcher committing a batch containing at most 10 operations and at most 1 second:
//
//	New(commitFn, WithMaxSize(10), WithTimeout(1 * time.Second))
func New(commitFn CommitFunc, opts ...Option) (*Batcher, error) {
	log := zerolog.New(os.Stdout).
		With().
		Timestamp().
		Logger().
		Level(zerolog.InfoLevel)

	batcher := &Batcher{
		commitFn: commitFn,
		maxSize:  UnlimitedSize,
		timeout:  NoTimeout,
		in:       make(chan *Operation, DefaultBuffer),
		log:      &log,
		closing:  0,
	}

	for _, opt := range opts {
		opt(batcher)
	}

	if batcher.commitFn == nil {
		return nil, ErrNilCommitFunc
	}

	if batcher.maxSize < 0 {
		return nil, ErrNegativeMaxSize
	}

	if batcher.timeout < 0 {
		return nil, ErrNegativeTimeout
	}

	if batcher.maxSize == UnlimitedSize && batcher.timeout == NoTimeout {
		return nil, ErrBadConditions
	}

	return batcher, nil
}

// Send creates a new operation and sends it to the batcher in a blocking
// fashion. If the provided context expires before the batcher receives the
// operation, Send returns the context's error.
func (b *Batcher) Send(ctx context.Context, v any) (*Operation, error) {
	if atomic.LoadUint32(&b.closing) == 1 {
		return nil, ErrClosed
	}

	op := newOperation(v)
	select {
	case b.in <- op:
		return op, nil
	case <-ctx.Done():
		b.logLost()

		return nil, ErrMissedValue
	}
}

// Batch receives operations from the batcher, calling the commit function
// whenever max size is reached or a timeout occurs.
// When the provided context expires, the batching process is interrupted and
// the function returns after a final call to the commit function.
// While shutting down, the send method would return error in a non-blocking fashion.
func (b *Batcher) Batch(ctx context.Context) {
	var tch <-chan time.Time

	out := b.makeBatch()
	timer := b.makeTimer()

	if timer != nil {
		tch = timer.C
	}

	b.logStart()

	for {
		var commit, done bool
		select {
		case op := <-b.in:
			out = append(out, op)
			commit = len(out) == b.maxSize
		case <-tch:
			commit = true
		case <-ctx.Done():
			atomic.StoreUint32(&b.closing, 1)
			close(b.in)

			// drain remaining operations
			for op := range b.in {
				out = append(out, op)
			}

			commit = len(out) > 0
			done = true
		}

		if commit {
			if len(out) > 0 {
				b.logCommit(len(out))
				b.commitFn(ctx, out)
			}

			out = out[:0]

			if timer != nil {
				timer.Reset(b.timeout)
				tch = timer.C
			} else {
				tch = nil
			}
		}

		if done {
			b.logEnd()

			break
		}
	}
}

func (b *Batcher) makeBatch() Batch {
	if b.maxSize != UnlimitedSize {
		return make(Batch, 0, b.maxSize)
	}

	return Batch{}
}

func (b *Batcher) makeTimer() *time.Timer {
	if b.timeout != NoTimeout {
		return time.NewTimer(b.timeout)
	}

	return nil
}

func (b *Batcher) logLost() {
	b.log.Error().
		Int("buffer_capacity", cap(b.in)).
		Int("buffer_size", len(b.in)).
		Msg("batcher: value lost")
}

func (b *Batcher) logStart() {
	b.log.Info().
		Int("maxSize", b.maxSize).
		Dur("timeout(ms)", b.timeout).
		Int("buffer_capacity", cap(b.in)).
		Int("buffer_size", len(b.in)).
		Msg("batcher: started")
}

func (b *Batcher) logEnd() {
	b.log.Info().
		Int("buffer_capacity", cap(b.in)).
		Int("buffer_size", len(b.in)).
		Msg("batcher: stopped gracefully")
}

func (b *Batcher) logCommit(size int) {
	b.log.Info().
		Int("batch_size", size).
		Msg("batcher: committing batch")
}
