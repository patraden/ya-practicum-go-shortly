package main

import (
	"net/http"

	"github.com/patraden/ya-practicum-go-shortly/internal/app/config"
	"github.com/patraden/ya-practicum-go-shortly/internal/app/handler"
	"github.com/patraden/ya-practicum-go-shortly/internal/app/logger"
	"github.com/patraden/ya-practicum-go-shortly/internal/app/service"
	"github.com/rs/zerolog"
)

func main() {
	log := logger.NewLogger(zerolog.InfoLevel).GetLogger()

	cfg := config.LoadConfig()
	service := service.NewInFileShortenerService(cfg, log)
	r := handler.NewRouter(service, cfg)

	err := http.ListenAndServe(cfg.ServerAddr, r)
	if err != nil {
		log.Fatal().Send()
	}
}
