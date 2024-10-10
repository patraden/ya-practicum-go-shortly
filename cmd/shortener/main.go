package main

import (
	"log"
	"net/http"

	"github.com/patraden/ya-practicum-go-shortly/internal/app/config"
	"github.com/patraden/ya-practicum-go-shortly/internal/app/handler"
	"github.com/patraden/ya-practicum-go-shortly/internal/app/service"
)

func main() {
	service := service.NewShortenerService()
	cfg := config.LoadConfig()
	r := handler.NewRouter(service, cfg)
	log.Fatal(http.ListenAndServe(cfg.ServerAddr, r))
}
