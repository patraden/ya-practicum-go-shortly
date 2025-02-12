package batcher

import "context"

type CommitFunc func(context.Context, Batch)
