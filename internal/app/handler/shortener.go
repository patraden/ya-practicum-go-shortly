package handler

import (
	"errors"
	"io"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/mailru/easyjson"
	"github.com/rs/zerolog"

	"github.com/patraden/ya-practicum-go-shortly/internal/app/config"
	"github.com/patraden/ya-practicum-go-shortly/internal/app/domain"
	e "github.com/patraden/ya-practicum-go-shortly/internal/app/domain/errors"
	"github.com/patraden/ya-practicum-go-shortly/internal/app/dto"
	"github.com/patraden/ya-practicum-go-shortly/internal/app/middleware"
	"github.com/patraden/ya-practicum-go-shortly/internal/app/service/shortener"
)

// ShortenerHandler provides HTTP request handling for URL shortening operations.
type ShortenerHandler struct {
	service shortener.URLShortener
	config  *config.Config
	log     *zerolog.Logger
}

// NewShortenerHandler creates a new instance of ShortenerHandler.
func NewShortenerHandler(service shortener.URLShortener, config *config.Config, log *zerolog.Logger) *ShortenerHandler {
	return &ShortenerHandler{
		service: service,
		config:  config,
		log:     log,
	}
}

// InsistentShortenerHandler creates a new instance of ShortenerHandler with InsistentShortener.
func InsistentShortenerHandler(
	service *shortener.InsistentShortener,
	config *config.Config,
	log *zerolog.Logger,
) *ShortenerHandler {
	return &ShortenerHandler{
		service: service,
		config:  config,
		log:     log,
	}
}

// RegisterRoutes register all handler routes within http router.
func (h *ShortenerHandler) RegisterRoutes(router chi.Router) {
	router.Group(func(r chi.Router) {
		r.Use(middleware.Authenticate(h.log, h.config))
		r.Get("/{shortURL}", h.HandleGetOriginalURL)
		r.Get("/api/user/urls", h.HandleGetUserURLs)
		r.Post("/api/shorten/batch", h.HandleBatchShortenURLJSON)
		r.Post("/api/shorten", h.HandleShortenURLJSON)
		r.Post("/", h.HandleShortenURL)
	})
}

// HandleGetOriginalURL handles requests to retrieve the original URL from a shortened slug.
func (h *ShortenerHandler) HandleGetOriginalURL(w http.ResponseWriter, r *http.Request) {
	slug := domain.Slug(chi.URLParam(r, "shortURL"))
	original, err := h.service.GetOriginalURL(r.Context(), slug)

	switch {
	case errors.Is(err, e.ErrSlugInvalid):
		http.Error(w, err.Error(), http.StatusBadRequest)

		return
	case errors.Is(err, e.ErrSlugNotFound):
		http.Error(w, err.Error(), http.StatusNotFound)

		return
	case errors.Is(err, e.ErrSlugDeleted):
		http.Error(w, err.Error(), http.StatusGone)

		return
	case errors.Is(err, e.ErrShortenerInternal) || err != nil:
		http.Error(w, err.Error(), http.StatusInternalServerError)

		return
	}

	w.Header().Add("Location", original.String())
	w.WriteHeader(http.StatusTemporaryRedirect)
}

// HandleGetUserURLs retrieves all shortened URLs for the requesting user.
func (h *ShortenerHandler) HandleGetUserURLs(w http.ResponseWriter, r *http.Request) {
	batch, err := h.service.GetUserURLs(r.Context())

	if errors.Is(err, e.ErrShortenerInternal) {
		http.Error(w, err.Error(), http.StatusInternalServerError)

		return
	}

	w.Header().Set(ContentType, ContentTypeJSON)

	if errors.Is(err, e.ErrUserNotFound) {
		w.WriteHeader(http.StatusNoContent)

		return
	}

	if _, err = easyjson.MarshalToWriter(batch, w); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)

		return
	}
}

// HandleShortenURL handles requests to shorten a given URL.
func (h *ShortenerHandler) HandleShortenURL(w http.ResponseWriter, r *http.Request) {
	b, err := io.ReadAll(r.Body)
	originalURL := domain.OriginalURL(string(b))

	if r.URL.Path != "/" || r.Body == http.NoBody || err != nil || !originalURL.IsValid() {
		http.Error(w, "bad request", http.StatusBadRequest)

		return
	}

	slug, err := h.service.ShortenURL(r.Context(), originalURL)
	if err != nil && !errors.Is(err, e.ErrOriginalExists) {
		http.Error(w, err.Error(), http.StatusInternalServerError)

		return
	}

	w.Header().Set(ContentType, ContentTypeText)

	if errors.Is(err, e.ErrOriginalExists) {
		w.WriteHeader(http.StatusConflict)
	} else {
		w.WriteHeader(http.StatusCreated)
	}

	if _, err = w.Write([]byte(slug.WithBaseURL(h.config.BaseURL))); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)

		return
	}
}

// HandleShortenURLJSON handles URL shortening requests with a JSON payload.
func (h *ShortenerHandler) HandleShortenURLJSON(w http.ResponseWriter, r *http.Request) {
	urlReq := dto.ShortenURLRequest{LongURL: ""}

	if err := easyjson.UnmarshalFromReader(r.Body, &urlReq); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)

		return
	}

	slug, err := h.service.ShortenURL(r.Context(), domain.OriginalURL(urlReq.LongURL))
	if err != nil && !errors.Is(err, e.ErrOriginalExists) {
		http.Error(w, err.Error(), http.StatusInternalServerError)

		return
	}

	urlResp := dto.ShortenedURLResponse{ShortURL: slug.WithBaseURL(h.config.BaseURL)}

	w.Header().Set(ContentType, ContentTypeJSON)

	if errors.Is(err, e.ErrOriginalExists) {
		w.WriteHeader(http.StatusConflict)
	} else {
		w.WriteHeader(http.StatusCreated)
	}

	if _, err = easyjson.MarshalToWriter(&urlResp, w); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)

		return
	}
}

// HandleBatchShortenURLJSON processes batch URL shortening requests.
func (h *ShortenerHandler) HandleBatchShortenURLJSON(w http.ResponseWriter, r *http.Request) {
	var urlReqs dto.OriginalURLBatch

	if err := easyjson.UnmarshalFromReader(r.Body, &urlReqs); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)

		return
	}

	batch, err := h.service.ShortenURLBatch(r.Context(), &urlReqs)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)

		return
	}

	result := make(dto.SlugBatch, len(*batch))
	for i, elem := range *batch {
		result[i] = dto.CorrelatedSlug{
			CorrelationID: elem.CorrelationID,
			Slug:          domain.Slug(elem.Slug.WithBaseURL(h.config.BaseURL)),
		}
	}

	w.Header().Set(ContentType, ContentTypeJSON)
	w.WriteHeader(http.StatusCreated)

	if _, err := easyjson.MarshalToWriter(result, w); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)

		return
	}
}
