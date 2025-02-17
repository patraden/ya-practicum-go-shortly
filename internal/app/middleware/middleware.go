package middleware

import (
	"compress/flate"
	"net/http"

	"github.com/go-chi/chi/v5/middleware"
)

const errLabel = "middleware"

// Compress returns middleware that compresses the response body.
func Compress() func(next http.Handler) http.Handler {
	return middleware.Compress(flate.DefaultCompression, "application/json", "text/plain")
}

// Recoverer returns middleware that recovers from panics and writes a 500 internal server error response.
func Recoverer() func(next http.Handler) http.Handler {
	return middleware.Recoverer
}

// StripSlashes returns middleware that removes trailing slashes from the request URL path.
func StripSlashes() func(next http.Handler) http.Handler {
	return middleware.StripSlashes
}
