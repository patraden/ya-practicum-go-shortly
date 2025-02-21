package batcher

import "context"

// CommitFunc is being called each time batcher completes a batch of operations according to its options.
type CommitFunc func(context.Context, Batch)
