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
	"github.com/patraden/ya-practicum-go-shortly/internal/app/service/shortener"
	"github.com/patraden/ya-practicum-go-shortly/internal/app/service/urlgenerator"
	"github.com/patraden/ya-practicum-go-shortly/internal/app/utils/postgres"
)

type Server struct {
	*http.Server
	config *config.Config
	repo   repository.URLRepository
	log    zerolog.Logger
}

func NewServer(
	config *config.Config,
	repo repository.URLRepository,
	gen urlgenerator.URLGenerator,
	log zerolog.Logger,
	db *postgres.Database,
) *Server {
	srv := shortener.NewInsistentShortener(repo, gen, config, log)
	shandler := handler.NewShortenerHandler(srv, config, log)
	phandler := handler.NewPingHandler(db, config, log)
	router := handler.NewRouter(shandler, phandler, log)

	return &Server{
		Server: &http.Server{
			Addr:              config.ServerAddr,
			Handler:           router,
			ReadHeaderTimeout: config.ServerReadHeaderTimeout,
			WriteTimeout:      config.ServerWriteTimeout,
			IdleTimeout:       config.ServerIdleTimeout,
		},
		config: config,
		repo:   repo,
		log:    log,
	}
}

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
	s.log.Info().Msg("Server started")
}

func (s *Server) WaitForShutdown(stopChan <-chan os.Signal) {
	<-stopChan
	s.log.Info().Msg("Shutdown signal received")

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
