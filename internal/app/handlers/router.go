package handlers

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/patraden/ya-practicum-go-shortly/internal/app/config"
)

func NewRouter(appConfig *config.Config) http.Handler {
	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Get("/{shortURL}", HandleLinkRepoGet(appConfig))
	r.Post("/", HandleLinkRepoPost(appConfig))
	r.NotFound(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "path not found", http.StatusBadRequest)
	}))

	return r
}
