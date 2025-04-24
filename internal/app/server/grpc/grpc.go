package grpc

import (
	"context"
	"net"

	"github.com/rs/zerolog"
	"google.golang.org/grpc"

	pb "github.com/patraden/ya-practicum-go-shortly/api/shortener/v1"
	"github.com/patraden/ya-practicum-go-shortly/internal/app/config"
	e "github.com/patraden/ya-practicum-go-shortly/internal/app/domain/errors"
	"github.com/patraden/ya-practicum-go-shortly/internal/app/handler"
	"github.com/patraden/ya-practicum-go-shortly/internal/app/middleware"
)

// Server represents the gRPC server for the URL shortening handler.
type Server struct {
	grpcServer *grpc.Server
	config     *config.Config
	handler    *handler.GRPCShortenerHandler
	log        *zerolog.Logger
}

// NewServer creates instance of Server.
func NewServer(config *config.Config, handler *handler.GRPCShortenerHandler, log *zerolog.Logger) *Server {
	return &Server{
		grpcServer: &grpc.Server{},
		config:     config,
		handler:    handler,
		log:        log,
	}
}

// Run starts the application server.
func (s *Server) Run() error {
	listen, err := net.Listen("tcp", s.config.ServerGRPCAddr)
	if err != nil {
		return err
	}

	intercepters := s.handler.Interceptors()

	intercepters = append(intercepters, middleware.WithLoggingInterceptor(s.log))
	s.grpcServer = grpc.NewServer(grpc.ChainUnaryInterceptor(intercepters...))

	pb.RegisterURLShortenerServiceServer(s.grpcServer, s.handler)

	return s.grpcServer.Serve(listen)
}

// Shutdown stops the application server.
func (s *Server) Shutdown(ctx context.Context) error {
	stopped := make(chan struct{})

	go func() {
		s.grpcServer.GracefulStop()
		close(stopped)
	}()

	select {
	case <-ctx.Done():
		s.grpcServer.Stop()
		return e.ErrServerShutdown
	case <-stopped:
		s.log.Info().Msg("gRPC server stopped gracefully")

		return nil
	}
}
