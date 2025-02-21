package server

import (
	"context"
	"errors"
	"net/http"
	"os"

	"github.com/rs/zerolog"

	"github.com/patraden/ya-practicum-go-shortly/internal/app/config"
	"github.com/patraden/ya-practicum-go-shortly/internal/app/handler"
	"github.com/patraden/ya-practicum-go-shortly/internal/app/memento"
	"github.com/patraden/ya-practicum-go-shortly/internal/app/repository"
	"github.com/patraden/ya-practicum-go-shortly/internal/app/service/remover"
	"github.com/patraden/ya-practicum-go-shortly/internal/app/service/shortener"
	"github.com/patraden/ya-practicum-go-shortly/internal/app/service/urlgenerator"
	"github.com/patraden/ya-practicum-go-shortly/internal/app/utils/postgres"
)

// Server represents the HTTP server for the URL shortening service.
type Server struct {
	*http.Server
	config        *config.Config
	repo          repository.URLRepository
	log           *zerolog.Logger
	remover       *remover.BatchRemover // this is really bad but I have no time :( now
	removerCancel func()
}

// NewServer creates a new Server instance with the provided configuration,
// repository, URL generator, logger, database, and batch remover.
func NewServer(
	config *config.Config,
	repo repository.URLRepository,
	gen urlgenerator.URLGenerator,
	log *zerolog.Logger,
	db *postgres.Database,
	remover *remover.BatchRemover,
) *Server {
	srv := shortener.NewInsistentShortener(repo, gen, config, log)
	shandler := handler.NewShortenerHandler(srv, config, log)
	phandler := handler.NewPingHandler(db, config, log)
	dhandler := handler.NewDeleteHandler(remover, log)
	router := handler.NewRouter(shandler, phandler, dhandler, log, config)

	return &Server{
		Server: &http.Server{
			Addr:              config.ServerAddr,
			Handler:           router,
			ReadHeaderTimeout: config.ServerReadHeaderTimeout,
			WriteTimeout:      config.ServerWriteTimeout,
			IdleTimeout:       config.ServerIdleTimeout,
		},
		config:        config,
		repo:          repo,
		remover:       remover,
		removerCancel: func() {},
		log:           log,
	}
}

// Start starts the HTTP server and loads the repository state from file if necessary.
// It also starts the batch remover service.
func (s *Server) Start() {
	if originator, ok := s.repo.(memento.Originator); ok && !s.config.ForceEmptyRepo {
		manager := memento.NewStateManager(s.config, originator, s.log)

		s.log.
			Info().
			Str("file_storage_path", s.config.FileStoragePath).
			Msg("restoring repository from file")
		s.loadRepositoryState(manager)
	}

	go func() {
		if err := s.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			s.log.Fatal().
				Err(err).
				Msg("Server failed")
		}
	}()

	ctx, cancel := context.WithCancel(context.Background())
	s.remover.Start(ctx)
	s.removerCancel = cancel

	s.log.Info().Msg("Server started")
}

// WaitForShutdown waits for a shutdown signal and handles graceful shutdown of the server and batch remover service.
func (s *Server) WaitForShutdown(stopChan <-chan os.Signal) {
	<-stopChan
	s.log.Info().Msg("Shutdown signal received")

	s.removerCancel()
	s.remover.Stop(context.Background())

	if originator, ok := s.repo.(memento.Originator); ok && !s.config.ForceEmptyRepo {
		manager := memento.NewStateManager(s.config, originator, s.log)
		s.log.
			Info().
			Str("file_storage_path", s.config.FileStoragePath).
			Msg("saving repository to file")
		s.saveRepositoryState(manager)
	}

	s.shutdownWithTimeout()
}

func (s *Server) shutdownWithTimeout() {
	ctx, cancel := context.WithTimeout(context.Background(), s.config.ServerShutTimeout)
	defer cancel()

	if err := s.Shutdown(ctx); err != nil {
		s.log.Error().
			Err(err).
			Msg("Server forced to shut down")
	} else {
		s.log.Info().
			Msg("Server shut down gracefully")
	}
}

func (s *Server) saveRepositoryState(manager *memento.StateManager) {
	if err := manager.StoreToFile(); err != nil {
		s.log.Error().
			Err(err).
			Str("file_storage_path", s.config.FileStoragePath).
			Msg("Failed to store state to file")
	}
}

func (s *Server) loadRepositoryState(manager *memento.StateManager) {
	if err := manager.RestoreFromFile(); err != nil {
		s.log.Error().
			Err(err).
			Str("file_storage_path", s.config.FileStoragePath).
			Msg("Failed to restore repository state from file")
	}
}
