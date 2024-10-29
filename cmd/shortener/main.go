package main

import (
	"context"
	"errors"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/patraden/ya-practicum-go-shortly/internal/app/config"
	"github.com/patraden/ya-practicum-go-shortly/internal/app/logger"
	"github.com/patraden/ya-practicum-go-shortly/internal/app/repository"
	"github.com/patraden/ya-practicum-go-shortly/internal/app/server"
	"github.com/patraden/ya-practicum-go-shortly/internal/app/urlgenerator"
	"github.com/rs/zerolog"
)

func main() {
	log := logger.NewLogger(zerolog.InfoLevel).GetLogger()
	config := config.LoadConfig()
	repo := repository.NewInMemoryURLRepository()
	gen := urlgenerator.NewRandURLGenerator(config.URLsize)
	srv := server.NewServer(repo, gen, config, log)
	manager := repository.NewStateManager(config, log)

	// Initialize repository with the help of repo state manager.
	memento, err := manager.LoadFromFile()
	if err != nil {
		log.Error().
			Err(err).
			Str("file_storage_path", config.FileStoragePath).
			Msg("failed to load state from file")
	} else {
		err := repo.RestoreMemento(memento)
		if err != nil {
			log.Error().
				Err(err).
				Str("file_storage_path", config.FileStoragePath).
				Msg("failed to restore repo from file")
		}
	}

	stopChan := make(chan os.Signal, 1)
	signal.Notify(stopChan, syscall.SIGINT, syscall.SIGTERM)

	// Start the server in a goroutine
	go func() {
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatal().
				Err(err).
				Msg("Server failed")
		}
	}()
	log.Info().Msg("Server started")

	// Wait for a termination signal.
	<-stopChan
	log.Info().Msg("Shutdown signal received")

	// Ensure we save the state on shutdown
	memento, err = repo.CreateMemento()
	if err != nil {
		log.Error().
			Err(err).
			Msg("failed to create repository state")
	} else {
		if err := manager.SaveToFile(memento); err != nil {
			log.Error().
				Err(err).
				Str("file_storage_path", config.FileStoragePath).
				Msg("failed to store state to file")
		}
	}

	ctx, cancel := context.WithTimeout(context.Background(), config.ServerShutTimeout)
	defer cancel()

	// Attempt to gracefully shut down the server.
	if err := srv.Shutdown(ctx); err != nil {
		log.Error().
			Err(err).
			Msg("Server forced to shut down")
	} else {
		log.Info().
			Msg("Server shut down gracefully")
	}
}
