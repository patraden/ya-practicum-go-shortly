package handler_test

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"github.com/patraden/ya-practicum-go-shortly/internal/app"
	"github.com/patraden/ya-practicum-go-shortly/internal/app/config"
)

func Example() {
	cfg := config.DefaultConfig()
	cfg.ForceEmptyRepo = true

	// get http client
	reqTimeout := 5 * time.Second
	client := http.DefaultClient
	ctx, cancel := context.WithTimeout(context.Background(), reqTimeout)

	defer cancel()

	// start application
	startTimeout := 2 * time.Second
	app := app.App(cfg, zerolog.Disabled)

	go func() { app.Run() }()

	time.Sleep(startTimeout)

	// send request
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

	if err = app.Stop(ctx); err != nil {
		log.Error().Err(err).Msg("failed stop app")
	}

	// Output:
	// 201 Created
}
