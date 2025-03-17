package handler

import "github.com/go-chi/chi/v5"

// Handler aux constants.
const (
	ContentType     = "Content-Type"
	ContentTypeText = "text/plain"
	ContentTypeJSON = "application/json"
)

// Handler can register its routes within router.
type Handler interface {
	RegisterRoutes(router chi.Router)
}
