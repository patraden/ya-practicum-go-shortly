package handler

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/rs/zerolog"

	"github.com/patraden/ya-practicum-go-shortly/internal/app/config"
	"github.com/patraden/ya-practicum-go-shortly/internal/app/utils/postgres"
)

// PingHandler handles database health checks.
type PingHandler struct {
	db     *postgres.Database
	config *config.Config
	log    *zerolog.Logger
}

// NewPingHandler creates and returns a new PingHandler.
func NewPingHandler(db *postgres.Database, config *config.Config, log *zerolog.Logger) *PingHandler {
	return &PingHandler{
		db:     db,
		config: config,
		log:    log,
	}
}

// RegisterRoutes register all handler routes within http router.
func (h *PingHandler) RegisterRoutes(router chi.Router) {
	router.Get("/ping", h.HandleDBPing)
}

// HandleDBPing verifies database connectivity and returns the status.
func (h *PingHandler) HandleDBPing(w http.ResponseWriter, r *http.Request) {
	if err := h.db.Ping(r.Context()); err != nil {
		h.log.
			Error().
			Err(err).
			Str("DSN", h.config.DatabaseDSN).
			Msg("database is not reachable")

		http.Error(w, "database is not reachable", http.StatusInternalServerError)

		return
	}

	w.WriteHeader(http.StatusOK)
}
