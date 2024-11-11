package handler_test

import (
	"bytes"
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/patraden/ya-practicum-go-shortly/internal/app/config"
	"github.com/patraden/ya-practicum-go-shortly/internal/app/domain"
	h "github.com/patraden/ya-practicum-go-shortly/internal/app/handler"
	"github.com/patraden/ya-practicum-go-shortly/internal/app/logger"
	"github.com/patraden/ya-practicum-go-shortly/internal/app/middleware"
	"github.com/patraden/ya-practicum-go-shortly/internal/app/repository"
	"github.com/patraden/ya-practicum-go-shortly/internal/app/service/shortener"
	"github.com/patraden/ya-practicum-go-shortly/internal/app/service/urlgenerator"
	"github.com/patraden/ya-practicum-go-shortly/internal/app/utils"
)

func setupEndPointTestRouter(repo repository.URLRepository) http.Handler {
	config := config.DefaultConfig()
	gen := urlgenerator.NewRandURLGenerator(config.URLsize)
	log := logger.NewLogger(zerolog.InfoLevel).GetLogger()
	srv := shortener.NewInsistentShortener(repo, gen, config, log)
	shandler := h.NewShortenerHandler(srv, config, log)

	return h.NewRouter(shandler, nil, log)
}

func TestEndpoints(t *testing.T) {
	t.Parallel()

	repo := repository.NewInMemoryURLRepository()
	err := repo.AddURLMapping(context.Background(), domain.NewURLMapping("shortURL", "https://ya.ru"))
	require.NoError(t, err)

	router := setupEndPointTestRouter(repo)

	tests := []struct {
		name   string
		method string
		path   string
		body   io.Reader
		want   int
	}{
		{"test 1", http.MethodDelete, "/", nil, http.StatusMethodNotAllowed},
		{"test 2", http.MethodPatch, "/", nil, http.StatusMethodNotAllowed},
		{"test 3", http.MethodPost, "/", strings.NewReader("https://ya.ru"), http.StatusCreated},
		{"test 4", http.MethodGet, "/shortURL", nil, http.StatusTemporaryRedirect},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			request := httptest.NewRequest(test.method, test.path, test.body)
			w := httptest.NewRecorder()
			router.ServeHTTP(w, request)

			result := w.Result()
			defer result.Body.Close()

			assert.Equal(t, test.want, result.StatusCode)

			if test.want == http.StatusTemporaryRedirect {
				assert.Equal(t, "https://ya.ru", result.Header.Get("Location"))
			}
		})
	}
}

func setupHandlerPost() http.HandlerFunc {
	config := config.DefaultConfig()
	repo := repository.NewInMemoryURLRepository()
	gen := urlgenerator.NewRandURLGenerator(config.URLsize)
	log := logger.NewLogger(zerolog.InfoLevel).GetLogger()
	srv := shortener.NewInsistentShortener(repo, gen, config, log)
	handler := h.NewShortenerHandler(srv, config, log)

	return handler.HandleShortenURL
}

func TestHandlePost(t *testing.T) {
	t.Parallel()

	handlePost := setupHandlerPost()

	tests := []struct {
		name        string
		path        string
		body        string
		contentType string
		wantStatus  int
	}{
		{"valid URL", "/", `https://ya.ru`, h.ContentTypeText, http.StatusCreated},
		{"empty body", "/", ``, h.ContentTypeText, http.StatusBadRequest},
		{"invalid URL", "/", `//ya.ru`, h.ContentTypeText, http.StatusBadRequest},
		{"invalid path", "/a/b/c", `https://ya.ru`, h.ContentTypeText, http.StatusBadRequest},
		{"double slash path", "//", `https://ya.ru`, h.ContentTypeText, http.StatusBadRequest},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			request := httptest.NewRequest(http.MethodPost, test.path, strings.NewReader(test.body))
			request.Header.Add(h.ContentType, test.contentType)

			w := httptest.NewRecorder()
			handlePost(w, request)

			result := w.Result()
			defer result.Body.Close()

			assert.Equal(t, test.wantStatus, result.StatusCode)
		})
	}
}

func setupHandleShortenURLJSON() http.HandlerFunc {
	config := config.DefaultConfig()
	repo := repository.NewInMemoryURLRepository()
	gen := urlgenerator.NewRandURLGenerator(config.URLsize)
	log := logger.NewLogger(zerolog.InfoLevel).GetLogger()
	srv := shortener.NewInsistentShortener(repo, gen, config, log)

	return h.NewShortenerHandler(srv, config, log).HandleShortenURLJSON
}

func TestHandleShortenURLJSON(t *testing.T) {
	t.Parallel()

	handlePostJSON := setupHandleShortenURLJSON()

	tests := []struct {
		name        string
		body        string
		contentType string
		wantStatus  int
	}{
		{"valid JSON URL", `{"url":"https://practicum.yandex.ru"}`, h.ContentTypeJSON, http.StatusCreated},
		{"malformed JSON", `{"url:"https://practicum.yandex.ru"}`, h.ContentTypeJSON, http.StatusBadRequest},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			request := httptest.NewRequest(http.MethodPost, `/api/shorten`, strings.NewReader(test.body))
			request.Header.Add(h.ContentType, test.contentType)

			w := httptest.NewRecorder()
			handlePostJSON(w, request)

			result := w.Result()
			defer result.Body.Close()

			assert.Equal(t, test.wantStatus, result.StatusCode)
		})
	}
}

func TestHandleGetOriginalURL(t *testing.T) {
	t.Parallel()

	config := config.DefaultConfig()
	log := logger.NewLogger(zerolog.InfoLevel).GetLogger()
	repo := repository.NewInMemoryURLRepository()
	gen := urlgenerator.NewRandURLGenerator(config.URLsize)
	srv := shortener.NewInsistentShortener(repo, gen, config, log)
	originalURL := domain.OriginalURL(`https://ya.ru`)
	baseURL := `http://localhost:8080/`

	link, err := srv.ShortenURL(context.Background(), originalURL)
	require.NoError(t, err)

	shandler := h.NewShortenerHandler(srv, config, log)

	router := h.NewRouter(shandler, nil, log)

	tests := []struct {
		name         string
		path         string
		wantStatus   int
		wantLocation string
	}{
		{"valid short URL", link.Slug.WithBaseURL(baseURL).String(), http.StatusTemporaryRedirect, string(originalURL)},
		{"invalid short URL", baseURL + "qwerty", http.StatusBadRequest, ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			request := httptest.NewRequest(http.MethodGet, tt.path, nil)
			w := httptest.NewRecorder()
			router.ServeHTTP(w, request)

			result := w.Result()
			defer result.Body.Close()

			assert.Equal(t, tt.wantStatus, result.StatusCode)

			if tt.wantStatus == http.StatusTemporaryRedirect {
				assert.Equal(t, tt.wantLocation, result.Header.Get("Location"))
			}
		})
	}
}

func setupHandleShortenURLCompression() http.Handler {
	config := config.DefaultConfig()
	repo := repository.NewInMemoryURLRepository()
	gen := urlgenerator.NewRandURLGenerator(config.URLsize)
	log := logger.NewLogger(zerolog.InfoLevel).GetLogger()
	srv := shortener.NewInsistentShortener(repo, gen, config, log)
	handler := http.HandlerFunc(h.NewShortenerHandler(srv, config, log).HandleShortenURL)

	return middleware.Decompress()(middleware.Compress()(handler))
}

func TestHandleShortenURLCompression(t *testing.T) {
	t.Parallel()

	handler := setupHandleShortenURLCompression()

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
