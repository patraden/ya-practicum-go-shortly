package handler

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/patraden/ya-practicum-go-shortly/internal/app/config"
	"github.com/patraden/ya-practicum-go-shortly/internal/app/middleware"
	"github.com/patraden/ya-practicum-go-shortly/internal/app/service"
)

func NewRouter(service *service.URLShortener, config *config.Config) http.Handler {
	handler := NewHandler(service, config)
	router := chi.NewRouter()

	router.Use(middleware.Recoverer())
	router.Use(middleware.StripSlashes())
	router.Use(middleware.Compress())
	router.Use(middleware.Decompress())
	router.Use(middleware.Logger())

	router.Get("/{shortURL}", handler.HandleGet)
	router.Post("/api/shorten", handler.HandlePostJSON)
	router.Post("/", handler.HandlePost)
	router.NotFound(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		http.Error(w, "path not found", http.StatusBadRequest)
	}))

	return router
}
