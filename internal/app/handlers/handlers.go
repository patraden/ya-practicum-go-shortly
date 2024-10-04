package handlers

import (
	"errors"
	"io"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/patraden/ya-practicum-go-shortly/internal/app/config"
	"github.com/patraden/ya-practicum-go-shortly/internal/app/helpers"
	"github.com/patraden/ya-practicum-go-shortly/internal/app/repository"
)

const (
	ContentType     = "Content-Type"
	ContentTypeText = "text/plain"
)

func HandleLinkRepoGet(appConfig *config.Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		shortURL := chi.URLParam(r, "shortURL")
		longURL, err := appConfig.Repo.ReStore(shortURL)

		if err != nil {
			if errors.Is(err, repository.ErrInternal) {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		w.Header().Add("Location", longURL)
		w.WriteHeader(http.StatusTemporaryRedirect)
	}
}

func HandleLinkRepoPost(appConfig *config.Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		b, err := io.ReadAll(r.Body)

		if r.URL.Path != "/" || r.Body == http.NoBody || err != nil || !helpers.IsURL(string(b)) {
			http.Error(w, "bad request", http.StatusBadRequest)
			return
		}

		shortURL, err := appConfig.Repo.Store(string(b))
		if err != nil {
			if errors.Is(err, repository.ErrInternal) {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		w.Header().Set(ContentType, ContentTypeText)
		w.WriteHeader(http.StatusCreated)
		w.Write([]byte(appConfig.BaseURL + shortURL))
	}
}
