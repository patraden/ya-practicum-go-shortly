package middleware

import (
	"net/http"
	"time"

	"github.com/patraden/ya-practicum-go-shortly/internal/app/logger"
)

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

func (r *loggingResponseWriter) Write(b []byte) (int, error) {
	size, err := r.ResponseWriter.Write(b)
	r.responseData.size += size
	return size, err
}

func (r *loggingResponseWriter) WriteHeader(statusCode int) {
	r.ResponseWriter.WriteHeader(statusCode)
	r.responseData.status = statusCode
}

func WithLogging(h http.Handler) http.Handler {
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

		logger.Log.Info().
			Str("uri", r.RequestURI).
			Str("method", r.Method).
			Int("status", responseData.status).
			Dur("duration", duration).
			Int("size", responseData.size).
			Msg("request details")
	}
	return http.HandlerFunc(logFn)
}
