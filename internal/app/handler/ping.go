package handler

import (
	"net/http"

	"github.com/rs/zerolog"

	"github.com/patraden/ya-practicum-go-shortly/internal/app/config"
	"github.com/patraden/ya-practicum-go-shortly/internal/app/utils/postgres"
)

type PingHandler struct {
	db     *postgres.Database
	config *config.Config
	log    *zerolog.Logger
}

func NewPingHandler(db *postgres.Database, config *config.Config, log *zerolog.Logger) *PingHandler {
	return &PingHandler{
		db:     db,
		config: config,
		log:    log,
	}
}

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
