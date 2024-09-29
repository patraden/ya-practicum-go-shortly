package handlers

import (
	"io"
	"net/http"
	"strings"

	"github.com/patraden/ya-practicum-go-shortly/internal/app/repository"
)

const (
	ContentType     = "Content-Type"
	ContentTypeText = "text/plain"
	ContentTypeJSON = "application/json"
)

func HandleLinkRepo(repo repository.LinkRepository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			HandleLinkRepoPost(repo).ServeHTTP(w, r)
		case http.MethodGet:
			HandleLinkRepoGet(repo).ServeHTTP(w, r)
		default:
			http.Error(w, "unknown method", http.StatusBadRequest)
		}
	}
}

func HandleLinkRepoGet(repo repository.LinkRepository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		shortURL := strings.TrimPrefix(r.URL.Path, "/")
		longURL, err := repo.ReStore(shortURL)
		if err != nil {
			// if err.Error() == "internal error" {
			// 	http.Error(w, err.Error(), http.StatusInternalServerError)
			// 	return
			// }
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

		if r.URL.Path != "/" || r.Body == http.NoBody || err != nil || r.Header.Get(ContentType) != ContentTypeText {
			http.Error(w, "bad request", http.StatusBadRequest)
			return
		}

		shortURL, err := repo.Store(string(b))
		if err != nil {
			// if err.Error() == "internal error" {
			// 	http.Error(w, err.Error(), http.StatusInternalServerError)
			// 	return
			// }
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		w.Header().Set(ContentType, ContentTypeText)
		w.WriteHeader(http.StatusCreated)
		w.Write([]byte("http://localhost:8080/" + shortURL))
	}
}
