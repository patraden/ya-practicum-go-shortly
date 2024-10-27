package middleware

import (
	"compress/flate"
	"net/http"

	"github.com/go-chi/chi/v5/middleware"
)

func Compress() func(next http.Handler) http.Handler {
	return middleware.Compress(flate.DefaultCompression, "application/json", "text/plain")
}

func Recoverer() func(next http.Handler) http.Handler {
	return middleware.Recoverer
}

func StripSlashes() func(next http.Handler) http.Handler {
	return middleware.StripSlashes
}
