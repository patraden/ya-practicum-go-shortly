package middleware

import (
	"net/http"
	"time"

	"github.com/rs/zerolog"

	e "github.com/patraden/ya-practicum-go-shortly/internal/app/domain/errors"
)

// Aux types.
type (
	responseData struct {
		status int
		size   int
	}

	loggingResponseWriter struct {
		http.ResponseWriter
		responseData *responseData
	}
)

// Write writes the response body and tracks the size of the response.
func (r *loggingResponseWriter) Write(b []byte) (int, error) {
	size, err := r.ResponseWriter.Write(b)
	r.responseData.size += size

	if err != nil {
		return size, e.Wrap("failed to write response", err, errLabel)
	}

	return size, nil
}

// WriteHeader sets the response status code and tracks it.
func (r *loggingResponseWriter) WriteHeader(statusCode int) {
	r.ResponseWriter.WriteHeader(statusCode)
	r.responseData.status = statusCode
}

// Logger is a middleware that logs details of the HTTP request and response.
// It logs the request URI, method, status code, duration, and response size.
func Logger(log *zerolog.Logger) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return WithLogging(next, log)
	}
}

// WithLogging is a middleware that logs details of the HTTP request and response.
// It logs the request URI, method, status code, duration, and response size.
func WithLogging(h http.Handler, log *zerolog.Logger) http.Handler {
	logFn := func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		responseData := &responseData{
			status: 0,
			size:   0,
		}
		lw := loggingResponseWriter{
			ResponseWriter: w,
			responseData:   responseData,
		}
		h.ServeHTTP(&lw, r)

		duration := time.Since(start)

		log.Info().
			Str("uri", r.RequestURI).
			Str("method", r.Method).
			Int("status", responseData.status).
			Dur("duration", duration).
			Int("size", responseData.size).
			Msg("request details")
	}

	return http.HandlerFunc(logFn)
}
