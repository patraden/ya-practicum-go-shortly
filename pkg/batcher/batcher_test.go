package batcher_test

import (
	"context"
	"strconv"
	"sync"
	"testing"
	"time"

	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/patraden/ya-practicum-go-shortly/pkg/batcher"
)

func TestNewBatcher(t *testing.T) {
	t.Parallel()

	for _, params := range []struct {
		name     string
		commitFn batcher.CommitFunc
		opts     []batcher.Option
		err      error
	}{
		{
			name:     "nil commit func",
			commitFn: nil,
			opts:     []batcher.Option{},
			err:      batcher.ErrNilCommitFunc,
		},
		{
			name:     "negative max size",
			opts:     []batcher.Option{batcher.WithMaxSize(-1)},
			commitFn: func(_ context.Context, _ batcher.Batch) {},
			err:      batcher.ErrNegativeMaxSize,
		},
		{
			name:     "negative timeout",
			opts:     []batcher.Option{batcher.WithTimeout(-1 * time.Second)},
			commitFn: func(_ context.Context, _ batcher.Batch) {},
			err:      batcher.ErrNegativeTimeout,
		},
		{
			name:     "unlimited size with no timeout",
			commitFn: func(_ context.Context, _ batcher.Batch) {},
			opts: []batcher.Option{
				batcher.WithMaxSize(batcher.UnlimitedSize),
				batcher.WithTimeout(batcher.NoTimeout),
			},
			err: batcher.ErrBadConditions,
		},
		{
			name:     "unlimited size with no timeout (no options)",
			commitFn: func(_ context.Context, _ batcher.Batch) {},
			opts:     []batcher.Option{},
			err:      batcher.ErrBadConditions,
		},
	} {
		t.Run(params.name, func(t *testing.T) {
			t.Parallel()

			_, err := batcher.New(params.commitFn, params.opts...)
			require.ErrorIs(t, err, params.err)
		})
	}
}

func TestBatcherSend(t *testing.T) {
	t.Parallel()

	t.Run("send value", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()

		commitFn := func(_ context.Context, _ batcher.Batch) {}
		b, err := batcher.New(commitFn, batcher.WithBufferSize(1), batcher.WithMaxSize(1))
		require.NoError(t, err)

		_, err = b.Send(ctx, []byte("value"))
		require.NoError(t, err)
	})

	t.Run("send value error", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
		defer cancel()

		commitFn := func(_ context.Context, _ batcher.Batch) {}

		b, err := batcher.New(commitFn, batcher.WithBufferSize(1), batcher.WithMaxSize(1))
		require.NoError(t, err)

		_, err = b.Send(ctx, []byte("value1"))
		require.NoError(t, err)

		_, err = b.Send(ctx, []byte("value2"))
		require.ErrorIs(t, err, batcher.ErrMissedValue)
	})
}

func TestBatcherBatchBySize(t *testing.T) {
	t.Parallel()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	countedTotalSize := 0
	maxSize := 10

	btr, err := batcher.New(
		func(_ context.Context, batch batcher.Batch) {
			if len(batch) == 0 {
				t.Error("unfilled batch committed")

				return
			}

			batch.SetError(nil)
			countedTotalSize += len(batch)
		},
		batcher.WithBufferSize(maxSize),
		batcher.WithMaxSize(maxSize),
	)

	require.NoError(t, err)

	wgb := sync.WaitGroup{}
	wgb.Add(1)

	go func() {
		defer wgb.Done()
		btr.Batch(ctx)
	}()

	for i := range maxSize * 10 {
		str := strconv.Itoa(i)
		val := []byte(str)

		_, err = btr.Send(ctx, val)
		require.NoError(t, err)
	}

	// Cancel the context to check that the batcher commits latent operations.
	cancel()
	wgb.Wait()

	assert.Equal(t, maxSize*10, countedTotalSize)
}

func TestBatcherBatchByTimeout(t *testing.T) {
	t.Parallel()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	countedTotalSize := 0
	timeout := 150 * time.Millisecond

	btr, err := batcher.New(
		func(_ context.Context, batch batcher.Batch) {
			if len(batch) == 0 {
				t.Error("unfilled batch committed")

				return
			}

			batch.SetError(nil)
			countedTotalSize += len(batch)

			time.Sleep(100 * time.Millisecond)
		},
		batcher.WithBufferSize(1),
		batcher.WithTimeout(timeout),
	)
	require.NoError(t, err)

	wgb := sync.WaitGroup{}
	wgb.Add(1)

	go func() {
		defer wgb.Done()
		btr.Batch(ctx)
	}()

	for i := range 100 {
		str := strconv.Itoa(i)
		val := []byte(str)

		_, err = btr.Send(ctx, val)
		require.NoError(t, err)

		time.Sleep(10 * time.Millisecond)
	}

	// Cancel the context to check that the batcher commits latent operations.
	cancel()
	wgb.Wait()

	assert.Equal(t, 100, countedTotalSize)
}

func setupBenchmarkBatchers(maxSize int, log zerolog.Logger) (*batcher.Batcher, *batcher.Batcher) {
	blow, err := batcher.New(
		func(_ context.Context, batch batcher.Batch) { batch.SetError(nil) },
		batcher.WithBufferSize(1),
		batcher.WithMaxSize(maxSize),
		batcher.WithTimeout(10*time.Millisecond),
		batcher.WithLogger(&log),
	)
	if err != nil {
		return nil, nil
	}

	bmax, err := batcher.New(
		func(_ context.Context, batch batcher.Batch) { batch.SetError(nil) },
		batcher.WithBufferSize(maxSize),
		batcher.WithTimeout(10*time.Millisecond),
		batcher.WithMaxSize(maxSize),
		batcher.WithLogger(&log),
	)
	if err != nil {
		return nil, nil
	}

	return blow, bmax
}

func BenchmarkBatcher(b *testing.B) {
	const maxSize = 100
	var wgb sync.WaitGroup

	log := zerolog.Nop()
	ctx, cancel := context.WithCancel(context.Background())

	defer cancel()

	blow, bmax := setupBenchmarkBatchers(maxSize, log)
	if blow == nil || bmax == nil {
		b.Fatal()
	}

	wgb.Add(2)

	go func() {
		defer wgb.Done()
		blow.Batch(ctx)
	}()

	go func() {
		defer wgb.Done()
		bmax.Batch(ctx)
	}()

	b.ResetTimer()
	b.Run("low capacity buffer", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_, err := blow.Send(ctx, []byte(strconv.Itoa(i)))
			if err != nil {
				b.Fatal(err)
			}
		}
	})

	b.Run("max capacity buffer", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_, err := bmax.Send(ctx, []byte(strconv.Itoa(i)))
			if err != nil {
				b.Fatal(err)
			}
		}
	})

	cancel()
	wgb.Wait()
}
