package http

import (
	"context"
	"net/http"

	"github.com/patraden/ya-practicum-go-shortly/internal/app/config"
)

// Server represents the HTTP server for the URL shortening service.
type Server struct {
	httpServer *http.Server
	config     *config.Config
}

// NewServer creates instance of Server.
func NewServer(config *config.Config, handler http.Handler) *Server {
	return &Server{
		httpServer: &http.Server{
			Addr:              config.ServerAddr,
			Handler:           handler,
			ReadHeaderTimeout: config.ServerReadHeaderTimeout,
			WriteTimeout:      config.ServerWriteTimeout,
			IdleTimeout:       config.ServerIdleTimeout,
		},
		config: config,
	}
}

// Run starts the application server.
func (s *Server) Run() error {
	if s.config.EnableHTTPS {
		return s.httpServer.ListenAndServeTLS(s.config.TLSCertPath, s.config.TLSKeyPath)
	}

	return s.httpServer.ListenAndServe()
}

// Shutdown stops the application server.
func (s *Server) Shutdown(ctx context.Context) error {
	return s.httpServer.Shutdown(ctx)
}
