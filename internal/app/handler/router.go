package handler

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/rs/zerolog"

	"github.com/patraden/ya-practicum-go-shortly/internal/app/config"
	"github.com/patraden/ya-practicum-go-shortly/internal/app/middleware"
)

func NewRouter(
	shandler *ShortenerHandler,
	phandler *PingHandler,
	dhandler *DeleteHandler,
	log *zerolog.Logger,
	config *config.Config,
) http.Handler {
	router := chi.NewRouter()

	// Apply common middleware to all routes
	router.Use(middleware.Recoverer())
	router.Use(middleware.StripSlashes())
	router.Use(middleware.Compress())
	router.Use(middleware.Decompress())

	// Routes without authentication
	router.Get("/ping", phandler.HandleDBPing)

	// Routes with authorization
	router.Group(func(r chi.Router) {
		r.Use(middleware.Authorize(log, config))
		r.Use(middleware.Logger(log))
		r.Delete("/api/user/urls", dhandler.HandleDelUserURLs)
	})

	// Routes with authentication
	router.Group(func(r chi.Router) {
		r.Use(middleware.Authenticate(log, config))
		r.Use(middleware.Logger(log))
		r.Get("/{shortURL}", shandler.HandleGetOriginalURL)
		r.Get("/api/user/urls", shandler.HandleGetUserURLs)
		r.Post("/api/shorten/batch", shandler.HandleBatchShortenURLJSON)
		r.Post("/api/shorten", shandler.HandleShortenURLJSON)
		r.Post("/", shandler.HandleShortenURL)
	})

	// Custom handler for 404
	router.NotFound(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		http.Error(w, "path not found", http.StatusBadRequest)
	}))

	return router
}
