package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/rs/zerolog"

	"github.com/patraden/ya-practicum-go-shortly/internal/app/config"
	"github.com/patraden/ya-practicum-go-shortly/internal/app/logger"
	"github.com/patraden/ya-practicum-go-shortly/internal/app/repository"
	"github.com/patraden/ya-practicum-go-shortly/internal/app/server"
	"github.com/patraden/ya-practicum-go-shortly/internal/app/service/remover"
	"github.com/patraden/ya-practicum-go-shortly/internal/app/service/urlgenerator"
	"github.com/patraden/ya-practicum-go-shortly/internal/app/utils/postgres"
)

func main() {
	var repo repository.URLRepository

	log := logger.NewLogger(zerolog.InfoLevel).GetLogger()
	cfg := config.LoadConfig()
	database := postgres.NewDatabase(log, cfg.DatabaseDSN)
	ctx := context.Background()
	gen := urlgenerator.NewRandURLGenerator(cfg.URLsize)
	repo = repository.NewInMemoryURLRepository()

	if cfg.DatabaseDSN != `` {
		err := database.Init(ctx)
		if err != nil {
			log.Fatal().
				Err(err).
				Msg("failed to connect to db repo")
		}

		repo = repository.NewDBURLRepository(database.ConnPool, log)
		cfg.ForceEmptyRepo = true
	}

	remover, err := remover.NewBatchRemover(repo, log)
	if err != nil {
		log.Fatal().
			Err(err).
			Msg("failed to init remover service")
	}

	srv := server.NewServer(cfg, repo, gen, log, database, remover)
	stopChan := make(chan os.Signal, 1)

	signal.Notify(stopChan, syscall.SIGINT, syscall.SIGTERM)
	srv.Start()
	srv.WaitForShutdown(stopChan)
}
