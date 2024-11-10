package handler

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/rs/zerolog"

	"github.com/patraden/ya-practicum-go-shortly/internal/app/middleware"
)

func NewRouter(
	shandler *ShortenerHandler,
	phandler *PingHandler,
	log zerolog.Logger,
) http.Handler {
	router := chi.NewRouter()

	router.Use(middleware.Recoverer())
	router.Use(middleware.StripSlashes())
	router.Use(middleware.Compress())
	router.Use(middleware.Decompress())
	router.Use(middleware.Logger(log))

	router.Get("/ping", phandler.HandleDBPing)
	router.Get("/{shortURL}", shandler.HandleGetOriginalURL)
	router.Post("/api/shorten", shandler.HandleShortenURLJSON)
	router.Post("/", shandler.HandleShortenURL)
	router.NotFound(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		http.Error(w, "path not found", http.StatusBadRequest)
	}))

	return router
}
