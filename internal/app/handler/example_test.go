package handler_test

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"syscall"
	"time"

	"github.com/rs/zerolog"

	"github.com/patraden/ya-practicum-go-shortly/internal/app/config"
	"github.com/patraden/ya-practicum-go-shortly/internal/app/logger"
	"github.com/patraden/ya-practicum-go-shortly/internal/app/repository"
	"github.com/patraden/ya-practicum-go-shortly/internal/app/server"
	"github.com/patraden/ya-practicum-go-shortly/internal/app/service/remover"
	"github.com/patraden/ya-practicum-go-shortly/internal/app/service/urlgenerator"
	"github.com/patraden/ya-practicum-go-shortly/internal/app/utils/postgres"
)

func Example() {
	cfg := config.LoadConfig()
	// avoid loading url snapshot from disc.
	cfg.ForceEmptyRepo = true

	log := logger.NewLogger(zerolog.Disabled).GetLogger()
	database := postgres.NewDatabase(log, cfg.DatabaseDSN)
	gen := urlgenerator.NewRandURLGenerator(cfg.URLsize)
	repo := repository.NewInMemoryURLRepository()
	client := http.DefaultClient

	remover, err := remover.NewBatchRemover(repo, log)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to init remover service")
	}

	stopChan := make(chan os.Signal, 1)
	srv := server.NewServer(cfg, repo, gen, log, database, remover)
	srv.Start()

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		"http://localhost:8080/",
		strings.NewReader("https://example.com"),
	)
	if err != nil {
		log.Error().Err(err).Msg("failed to create POST request to shorten URL")
	}

	req.Header.Set("Content-Type", "text/plain")

	resp, err := client.Do(req)
	if err != nil {
		log.Error().Err(err).Msg("failed to send POST request to shorten URL")
	}

	defer resp.Body.Close()

	// Read and process the response body
	_, err = io.ReadAll(resp.Body)
	if err != nil {
		log.Error().Err(err).Msg("failed to read response body")
	}

	// Print the status code
	fmt.Println(resp.Status)

	// Gracefully shut down the server
	stopChan <- syscall.SIGINT
	srv.WaitForShutdown(stopChan)

	// Output:
	// 201 Created
}
