package handler

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/patraden/ya-practicum-go-shortly/internal/app/config"
	"github.com/patraden/ya-practicum-go-shortly/internal/app/middleware"
	"github.com/patraden/ya-practicum-go-shortly/internal/app/service"
)

func NewRouter(service service.URLShortener, config config.Config) http.Handler {
	h := NewHandler(service, config)
	r := chi.NewRouter()

	r.Use(middleware.Recoverer())
	r.Use(middleware.StripSlashes())
	r.Use(middleware.Compress())
	r.Use(middleware.Decompress())
	r.Use(middleware.Logger())

	r.Get("/{shortURL}", h.HandleGet)
	r.Post("/api/shorten", h.HandlePostJSON)
	r.Post("/", h.HandlePost)
	r.NotFound(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		http.Error(w, "path not found", http.StatusBadRequest)
	}))

	return r
}
