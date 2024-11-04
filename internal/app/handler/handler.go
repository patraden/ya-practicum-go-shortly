package handler

import (
	"errors"
	"io"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/mailru/easyjson"
	"github.com/patraden/ya-practicum-go-shortly/internal/app/config"
	"github.com/patraden/ya-practicum-go-shortly/internal/app/dto"
	e "github.com/patraden/ya-practicum-go-shortly/internal/app/errors"
	"github.com/patraden/ya-practicum-go-shortly/internal/app/service"
	"github.com/patraden/ya-practicum-go-shortly/internal/app/utils"
	"github.com/rs/zerolog"
)

const (
	ContentType     = "Content-Type"
	ContentTypeText = "text/plain"
	ContentTypeJSON = "application/json"
)

type Handler struct {
	service service.URLShortener
	config  *config.Config
	log     zerolog.Logger
}

func NewHandler(service service.URLShortener, config *config.Config, log zerolog.Logger) *Handler {
	return &Handler{
		service: service,
		config:  config,
		log:     log,
	}
}

func (h *Handler) HandleGet(w http.ResponseWriter, r *http.Request) {
	shortURL := chi.URLParam(r, "shortURL")
	longURL, err := h.service.GetOriginalURL(shortURL)

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

	w.Header().Add("Location", longURL)
	w.WriteHeader(http.StatusTemporaryRedirect)
}

func (h *Handler) HandlePost(w http.ResponseWriter, r *http.Request) {
	b, err := io.ReadAll(r.Body)

	if r.URL.Path != "/" || r.Body == http.NoBody || err != nil || !utils.IsURL(string(b)) {
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

	if _, err = w.Write([]byte(h.withBaseURL(shortURL))); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)

		return
	}
}

func (h *Handler) HandlePostJSON(w http.ResponseWriter, r *http.Request) {
	urlReq := dto.ShortenURLRequest{LongURL: ""}

	if err := easyjson.UnmarshalFromReader(r.Body, &urlReq); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)

		return
	}

	shortURL, err := h.service.ShortenURL(urlReq.LongURL)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)

		return
	}

	urlResp := dto.ShortenedURLResponse{ShortURL: h.withBaseURL(shortURL)}

	w.Header().Set(ContentType, ContentTypeJSON)
	w.WriteHeader(http.StatusCreated)

	if _, err = easyjson.MarshalToWriter(&urlResp, w); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)

		return
	}
}

func (h *Handler) withBaseURL(shortURL string) string {
	return h.config.BaseURL + shortURL

}
