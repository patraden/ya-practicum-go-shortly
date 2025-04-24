package server

import "context"

// AppServer common interface for HTTP and gRPC impementations.
type AppServer interface {
	Run() error
	Shutdown(ctx context.Context) error
}
