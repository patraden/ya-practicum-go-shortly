package handler

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/mailru/easyjson"
	"github.com/rs/zerolog"

	"github.com/patraden/ya-practicum-go-shortly/internal/app/config"
	"github.com/patraden/ya-practicum-go-shortly/internal/app/middleware"
	"github.com/patraden/ya-practicum-go-shortly/internal/app/service/statsprovider"
)

// StatsProviderHandler provides HTTP request handling for URL shortening statistics.
type StatsProviderHandler struct {
	service statsprovider.StatsProvider
	config  *config.Config
	log     *zerolog.Logger
}

// NewStatsProviderHandler creates new instance of StatsProviderHandler.
func NewStatsProviderHandler(
	service statsprovider.StatsProvider,
	config *config.Config,
	log *zerolog.Logger,
) *StatsProviderHandler {
	return &StatsProviderHandler{
		service: service,
		config:  config,
		log:     log,
	}
}

// RegisterRoutes register all handler routes within http router.
func (h *StatsProviderHandler) RegisterRoutes(router chi.Router) {
	router.Group(func(r chi.Router) {
		r.Use(middleware.SubnetMiddleware(h.log, h.config))
		r.Get("/api/internal/stats", h.HandleGetStats)
	})
}

// HandleGetStats handles requests to retrieve repository statistics.
func (h *StatsProviderHandler) HandleGetStats(w http.ResponseWriter, r *http.Request) {
	stats, err := h.service.GetStats(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)

		return
	}

	w.Header().Set(ContentType, ContentTypeJSON)

	if _, err = easyjson.MarshalToWriter(stats, w); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)

		return
	}
}
