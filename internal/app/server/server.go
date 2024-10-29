package server

import (
	"net/http"

	"github.com/patraden/ya-practicum-go-shortly/internal/app/config"
	"github.com/patraden/ya-practicum-go-shortly/internal/app/handler"
	"github.com/patraden/ya-practicum-go-shortly/internal/app/repository"
	"github.com/patraden/ya-practicum-go-shortly/internal/app/service"
	"github.com/patraden/ya-practicum-go-shortly/internal/app/urlgenerator"
)

func NewServer(
	repo repository.URLRepository,
	gen urlgenerator.URLGenerator,
	config *config.Config,
) *http.Server {
	service := service.NewShortenerService(repo, gen, config)
	router := handler.NewRouter(service, config)

	return &http.Server{
		Addr:    config.ServerAddr,
		Handler: router,
	}
}
