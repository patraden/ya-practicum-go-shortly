package handler

import (
	"compress/flate"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/patraden/ya-practicum-go-shortly/internal/app/config"
	cmiddleware "github.com/patraden/ya-practicum-go-shortly/internal/app/middleware"
	"github.com/patraden/ya-practicum-go-shortly/internal/app/service"
)

func NewRouter(service service.URLShortener, config config.Config) http.Handler {
	h := NewHandler(service, config)
	r := chi.NewRouter()

	r.Use(middleware.Recoverer)
	r.Use(middleware.StripSlashes)
	r.Use(middleware.Compress(flate.DefaultCompression, ContentTypeJSON, ContentTypeText))
	r.Use(cmiddleware.Decompress)
	r.Use(cmiddleware.WithLogging)

	r.Get("/{shortURL}", h.HandleGet)
	r.Post("/api/shorten", h.HandlePostJSON)
	r.Post("/", h.HandlePost)
	r.NotFound(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "path not found", http.StatusBadRequest)
	}))

	return r
}
