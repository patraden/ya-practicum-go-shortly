package handlers

import (
	"io"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/patraden/ya-practicum-go-shortly/internal/app/repository"
)

const (
	ContentType     = "Content-Type"
	ContentTypeText = "text/plain"
	ContentTypeJSON = "application/json"
)

func HandleLinkRepoGet(repo repository.LinkRepository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		shortURL := chi.URLParam(r, "shortURL")
		longURL, err := repo.ReStore(shortURL)

		if err != nil {
			if err.Error() == "internal error" {
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

func HandleLinkRepoPost(repo repository.LinkRepository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		b, err := io.ReadAll(r.Body)

		if r.URL.Path != "/" || r.Body == http.NoBody || err != nil {
			http.Error(w, "bad request", http.StatusBadRequest)
			return
		}

		shortURL, err := repo.Store(string(b))
		if err != nil {
			if err.Error() == "internal error" {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		w.Header().Set(ContentType, ContentTypeText)
		w.WriteHeader(http.StatusCreated)
		w.Write([]byte("http://" + r.Host + "/" + shortURL))
	}
}
