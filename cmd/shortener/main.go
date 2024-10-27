package main

import (
	"log"
	"net/http"

	"github.com/patraden/ya-practicum-go-shortly/internal/app/config"
	"github.com/patraden/ya-practicum-go-shortly/internal/app/handler"
	"github.com/patraden/ya-practicum-go-shortly/internal/app/logger"
	"github.com/patraden/ya-practicum-go-shortly/internal/app/service"
	"github.com/rs/zerolog"
)

func main() {
	logger.NewLogger(zerolog.InfoLevel).GetLogger()

	cfg := config.LoadConfig()
	service := service.NewInMemoryShortenerService(cfg)
	r := handler.NewRouter(service, cfg)

	err := http.ListenAndServe(cfg.ServerAddr, r)
	if err != nil {
		log.Fatal(err)
	}
}
