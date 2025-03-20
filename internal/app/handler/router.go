package handler

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/rs/zerolog"

	"github.com/patraden/ya-practicum-go-shortly/internal/app/middleware"
)

// HTTP router.
func NewRouter(log *zerolog.Logger, handlers ...Handler) http.Handler {
	router := chi.NewRouter()

	// Apply common middleware to all routes
	router.Use(middleware.Recoverer())
	router.Use(middleware.StripSlashes())
	router.Use(middleware.Compress())
	router.Use(middleware.Decompress())
	router.Use(middleware.Logger(log))

	// Register handlers
	for _, h := range handlers {
		h.RegisterRoutes(router)
	}

	// Custom handler for 404
	router.NotFound(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		http.Error(w, "path not found", http.StatusBadRequest)
	}))

	return router
}
