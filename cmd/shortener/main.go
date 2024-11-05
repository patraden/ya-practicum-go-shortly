package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/rs/zerolog"

	"github.com/patraden/ya-practicum-go-shortly/internal/app/config"
	"github.com/patraden/ya-practicum-go-shortly/internal/app/logger"
	"github.com/patraden/ya-practicum-go-shortly/internal/app/repository"
	"github.com/patraden/ya-practicum-go-shortly/internal/app/server"
	"github.com/patraden/ya-practicum-go-shortly/internal/app/urlgenerator"
)

func main() {
	log := logger.NewLogger(zerolog.InfoLevel).GetLogger()
	cfg := config.LoadConfig()
	repo := repository.NewInMemoryURLRepository()
	mngr := repository.NewStateManager(cfg, log)
	gen := urlgenerator.NewRandURLGenerator(cfg.URLsize)

	srv := server.NewServer(cfg, repo, mngr, gen, log)
	stopChan := make(chan os.Signal, 1)

	signal.Notify(stopChan, syscall.SIGINT, syscall.SIGTERM)
	srv.Start()
	srv.WaitForShutdown(stopChan)
}
