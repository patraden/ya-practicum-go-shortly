package batcher

import "context"

// Operation represents a single asynchronous Batcher operation that can be awaited
// for completion. It holds the resulted error that occurred during the batch commit,
// and signals when the operation is done.
type Operation struct {
	Value any
	err   error
	done  chan struct{}
}

func newOperation(v any) *Operation {
	return &Operation{
		Value: v,
		err:   nil,
		done:  make(chan struct{}),
	}
}

// SetError signals an error relating to the operation.
func (o *Operation) SetError(err error) {
	o.err = err
	close(o.done)
}

// Wait blocks until the operation completes, returning the result or the error
// encountered. If the provided context expires before the operation is
// complete, Wait returns the context's error.
func (o *Operation) Wait(ctx context.Context) error {
	select {
	case <-o.done:
		return o.err
	case <-ctx.Done():
		return ErrOperationWait
	}
}

// IsDone checks if the operation has been completed.
func (o *Operation) IsDone() bool {
	select {
	case <-o.done:
		return true
	default:
		return false
	}
}
