package server

import (
	"net/http"

	"github.com/patraden/ya-practicum-go-shortly/internal/app/config"
	"github.com/patraden/ya-practicum-go-shortly/internal/app/handler"
	"github.com/patraden/ya-practicum-go-shortly/internal/app/repository"
	"github.com/patraden/ya-practicum-go-shortly/internal/app/service"
	"github.com/patraden/ya-practicum-go-shortly/internal/app/urlgenerator"
	"github.com/rs/zerolog"
)

func NewServer(
	repo repository.URLRepository,
	gen urlgenerator.URLGenerator,
	config *config.Config,
	log zerolog.Logger,
) *http.Server {

	if config == nil {
		log.Fatal().Msg("Config is not initialized")
	}

	service := service.NewShortenerService(repo, gen, config)
	router := handler.NewRouter(service, config)

	return &http.Server{
		Addr:    config.ServerAddr,
		Handler: router,
	}

}
