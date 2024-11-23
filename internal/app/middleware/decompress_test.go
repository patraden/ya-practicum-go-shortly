package middleware_test

import (
	"bytes"
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/patraden/ya-practicum-go-shortly/internal/app/config"
	"github.com/patraden/ya-practicum-go-shortly/internal/app/domain"
	"github.com/patraden/ya-practicum-go-shortly/internal/app/handler"
	"github.com/patraden/ya-practicum-go-shortly/internal/app/logger"
	"github.com/patraden/ya-practicum-go-shortly/internal/app/middleware"
	"github.com/patraden/ya-practicum-go-shortly/internal/app/repository"
	"github.com/patraden/ya-practicum-go-shortly/internal/app/service/shortener"
	"github.com/patraden/ya-practicum-go-shortly/internal/app/service/urlgenerator"
	"github.com/patraden/ya-practicum-go-shortly/internal/app/utils"
)

func setupHandleShortenURLCompression() http.Handler {
	config := config.DefaultConfig()
	repo := repository.NewInMemoryURLRepository()
	gen := urlgenerator.NewRandURLGenerator(config.URLsize)
	log := logger.NewLogger(zerolog.InfoLevel).GetLogger()
	srv := shortener.NewInsistentShortener(repo, gen, config, log)
	handler := http.HandlerFunc(handler.NewShortenerHandler(srv, config, log).HandleShortenURL)

	return middleware.Decompress()(middleware.Compress()(handler))
}

func TestHandleShortenURLCompression(t *testing.T) {
	t.Parallel()

	handler := setupHandleShortenURLCompression()
	userID := domain.NewUserID()

	tests := []struct {
		name            string
		contentEncoding string
		acceptEncoding  string
		status          int
		isURL           bool
	}{
		{"test 1", "gzip", "", http.StatusCreated, true},
		{"test 2", "deflate", "", http.StatusCreated, true},
		{"test 3", "gzip", "deflate", http.StatusCreated, true},
		{"test 4", "deflate", "gzip", http.StatusCreated, true},
		{"test 5", "", "", http.StatusCreated, true},
	}

	for _, test := range tests {
		for range 1 {
			t.Run(test.name, func(t *testing.T) {
				t.Parallel()

				url := utils.RandURL()
				data, err := utils.Compress([]byte(url), test.contentEncoding)
				require.NoError(t, err)

				r := httptest.NewRequest(http.MethodPost, "/", bytes.NewReader(data))
				r.Header.Set("Content-Encoding", test.contentEncoding)
				r.Header.Set("Accept-Encoding", test.acceptEncoding)

				r = r.Clone(context.WithValue(r.Context(), middleware.UserIDKey, userID))

				w := httptest.NewRecorder()
				handler.ServeHTTP(w, r)

				// Read the response and close the body
				result := w.Result()
				defer result.Body.Close() // Ensure the body is closed

				compressedURL, err := io.ReadAll(result.Body)
				require.NoError(t, err)

				shortURL, err := utils.Decompress(compressedURL, test.acceptEncoding)
				require.NoError(t, err)

				assert.Equal(t, test.status, result.StatusCode)
				assert.Equal(t, test.isURL, domain.OriginalURL(string(shortURL)).IsValid())
			})
		}
	}
}
