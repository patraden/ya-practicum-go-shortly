package batcher

import "errors"

// Batcher static errors.
var (
	ErrClosed          = errors.New("batcher: is closed")
	ErrMissedValue     = errors.New("batcher: missed value")
	ErrOperationWait   = errors.New("batcher: operation wait time exceeded")
	ErrNilCommitFunc   = errors.New("batcher: nil commit func")
	ErrNegativeMaxSize = errors.New("batcher: negative max size")
	ErrNegativeTimeout = errors.New("batcher: negative timeout")
	ErrBadConditions   = errors.New("batcher: unlimited size with no timeout")
)
