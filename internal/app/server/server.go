package server

import (
	"context"
	"errors"
	"net/http"
	"os"

	"github.com/patraden/ya-practicum-go-shortly/internal/app/config"
	"github.com/patraden/ya-practicum-go-shortly/internal/app/handler"
	"github.com/patraden/ya-practicum-go-shortly/internal/app/repository"
	"github.com/patraden/ya-practicum-go-shortly/internal/app/service"
	"github.com/patraden/ya-practicum-go-shortly/internal/app/urlgenerator"
	"github.com/rs/zerolog"
)

type Server struct {
	*http.Server
	config  *config.Config
	repo    repository.URLRepository
	manager *repository.StateManager
	log     zerolog.Logger
}

func NewServer(
	config *config.Config,
	repo repository.URLRepository,
	manager *repository.StateManager,
	gen urlgenerator.URLGenerator,
	log zerolog.Logger,
) *Server {
	service := service.NewShortenerService(repo, gen, config)
	router := handler.NewRouter(service, config, log)

	return &Server{
		Server: &http.Server{
			Addr:              config.ServerAddr,
			Handler:           router,
			ReadHeaderTimeout: config.ServerReadHeaderTimeout,
			WriteTimeout:      config.ServerWriteTimeout,
			IdleTimeout:       config.ServerIdleTimeout,
		},
		config:  config,
		repo:    repo,
		manager: manager,
		log:     log,
	}
}

func (s *Server) Start() {
	if !s.config.ForceEmptyRepo {
		s.log.
			Info().
			Str("file_storage_path", s.config.FileStoragePath).
			Msg("restoring repository from file")
		s.loadRepositoryState()
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

	if !s.config.ForceEmptyRepo {
		s.log.
			Info().
			Str("file_storage_path", s.config.FileStoragePath).
			Msg("saving repository to file")
		s.saveRepositoryState()
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

func (s *Server) saveRepositoryState() {
	memento, err := s.repo.CreateMemento()
	if err != nil {
		s.log.Error().
			Err(err).
			Msg("Failed to create repository state")

		return
	}

	if err := s.manager.SaveToFile(memento); err != nil {
		s.log.Error().
			Err(err).
			Str("file_storage_path", s.config.FileStoragePath).
			Msg("Failed to store state to file")
	}
}

func (s *Server) loadRepositoryState() {
	memento, err := s.manager.LoadFromFile()
	if err != nil {
		s.log.Error().
			Err(err).
			Str("file_storage_path", s.config.FileStoragePath).
			Msg("Failed to load state from file")

		return
	}

	if err := s.repo.RestoreMemento(memento); err != nil {
		s.log.Error().
			Err(err).
			Str("file_storage_path", s.config.FileStoragePath).
			Msg("Failed to restore repository state from file")
	}
}
