package main

import (
	"log"
	"net/http"

	"github.com/patraden/ya-practicum-go-shortly/internal/app/config"
	"github.com/patraden/ya-practicum-go-shortly/internal/app/handlers"
	"github.com/patraden/ya-practicum-go-shortly/internal/app/repository"
)

func main() {
	repo := repository.NewBasicLinkRepository()
	appConfig := config.LoadConfig(repo)
	r := handlers.NewRouter(appConfig)
	log.Fatal(http.ListenAndServe(appConfig.ServerAddr, r))
}
