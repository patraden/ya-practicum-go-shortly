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
	"github.com/patraden/ya-practicum-go-shortly/internal/app/dto"
	e "github.com/patraden/ya-practicum-go-shortly/internal/app/errors"
	"github.com/patraden/ya-practicum-go-shortly/internal/app/service/shortener"
)

const (
	ContentType     = "Content-Type"
	ContentTypeText = "text/plain"
	ContentTypeJSON = "application/json"
)

type ShortenerHandler struct {
	service shortener.URLShortener
	config  *config.Config
	log     zerolog.Logger
}

func NewShortenerHandler(service shortener.URLShortener, config *config.Config, log zerolog.Logger) *ShortenerHandler {
	return &ShortenerHandler{
		service: service,
		config:  config,
		log:     log,
	}
}

func (h *ShortenerHandler) HandleGetOriginalURL(w http.ResponseWriter, r *http.Request) {
	slug := domain.Slug(chi.URLParam(r, "shortURL"))
	link, err := h.service.GetOriginalURL(r.Context(), slug)

	switch {
	case errors.Is(err, e.ErrServiceInvalid):
		http.Error(w, err.Error(), http.StatusBadRequest)

		return
	case errors.Is(err, e.ErrRepoNotFound):
		http.Error(w, err.Error(), http.StatusNotFound)

		return
	case errors.Is(err, e.ErrServiceInternal) || err != nil:
		http.Error(w, err.Error(), http.StatusInternalServerError)

		return
	}

	w.Header().Add("Location", string(link.OriginalURL))
	w.WriteHeader(http.StatusTemporaryRedirect)
}

func (h *ShortenerHandler) HandleShortenURL(w http.ResponseWriter, r *http.Request) {
	b, err := io.ReadAll(r.Body)
	originalURL := domain.OriginalURL(string(b))

	if r.URL.Path != "/" || r.Body == http.NoBody || err != nil || !originalURL.IsValid() {
		http.Error(w, "bad request", http.StatusBadRequest)

		return
	}

	link, err := h.service.ShortenURL(r.Context(), originalURL)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)

		return
	}

	w.Header().Set(ContentType, ContentTypeText)
	w.WriteHeader(http.StatusCreated)

	if _, err = w.Write([]byte(link.Slug.WithBaseURL(h.config.BaseURL).String())); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)

		return
	}
}

func (h *ShortenerHandler) HandleShortenURLJSON(w http.ResponseWriter, r *http.Request) {
	urlReq := dto.ShortenURLRequest{LongURL: ""}

	if err := easyjson.UnmarshalFromReader(r.Body, &urlReq); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)

		return
	}

	link, err := h.service.ShortenURL(r.Context(), domain.OriginalURL(urlReq.LongURL))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)

		return
	}

	urlResp := dto.ShortenedURLResponse{ShortURL: link.Slug.WithBaseURL(h.config.BaseURL).String()}

	w.Header().Set(ContentType, ContentTypeJSON)
	w.WriteHeader(http.StatusCreated)

	if _, err = easyjson.MarshalToWriter(&urlResp, w); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)

		return
	}
}

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

	w.Header().Set(ContentType, ContentTypeJSON)
	w.WriteHeader(http.StatusCreated)

	if _, err := easyjson.MarshalToWriter(batch, w); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)

		return
	}
}
