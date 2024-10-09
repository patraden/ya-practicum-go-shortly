package handler

import (
	"errors"
	"io"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/patraden/ya-practicum-go-shortly/internal/app/config"
	e "github.com/patraden/ya-practicum-go-shortly/internal/app/errors"
	"github.com/patraden/ya-practicum-go-shortly/internal/app/helpers"
	"github.com/patraden/ya-practicum-go-shortly/internal/app/service"
)

const (
	ContentType     = "Content-Type"
	ContentTypeText = "text/plain"
)

type Handler struct {
	service service.URLShortener
	config  config.Config
}

func NewHandler(service service.URLShortener, config config.Config) *Handler {
	return &Handler{
		service: service,
		config:  config,
	}
}

func (h *Handler) HandleGet(w http.ResponseWriter, r *http.Request) {
	shortURL := chi.URLParam(r, "shortURL")
	longURL, err := h.service.GetOriginalURL(shortURL)

	if errors.Is(err, e.ErrInvalid) || errors.Is(err, e.ErrNotFound) {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if errors.Is(err, e.ErrInternal) || err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Add("Location", longURL)
	w.WriteHeader(http.StatusTemporaryRedirect)
}

func (h *Handler) HandlePost(w http.ResponseWriter, r *http.Request) {
	b, err := io.ReadAll(r.Body)

	if r.URL.Path != "/" || r.Body == http.NoBody || err != nil || !helpers.IsURL(string(b)) {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}

	shortURL, err := h.service.ShortenURL(string(b))

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set(ContentType, ContentTypeText)
	w.WriteHeader(http.StatusCreated)
	_, err = w.Write([]byte(h.config.BaseURL + shortURL))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
