package batcher

import "context"

type Operation struct {
	Value []byte
	err   error
	done  chan struct{}
}

func newOperation(v []byte) *Operation {
	return &Operation{
		Value: v,
		err:   nil,
		done:  make(chan struct{}),
	}
}

func (o *Operation) SetError(err error) {
	o.err = err
	close(o.done)
}

func (o *Operation) Wait(ctx context.Context) error {
	select {
	case <-o.done:
		return o.err
	case <-ctx.Done():
		return ErrOperationWait
	}
}

func (o *Operation) IsDone() bool {
	select {
	case <-o.done:
		return true
	default:
		return false
	}
}
