package handlers

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/patraden/ya-practicum-go-shortly/internal/app/repository"
)

func NewRouter(repo repository.LinkRepository) http.Handler {
	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Get("/{shortURL}", HandleLinkRepoGet(repo))
	r.Post("/", HandleLinkRepoPost(repo))
	r.NotFound(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "not found", http.StatusBadRequest)
	}))

	return r
}
