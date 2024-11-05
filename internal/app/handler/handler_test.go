package handler_test

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/patraden/ya-practicum-go-shortly/internal/app/config"
	h "github.com/patraden/ya-practicum-go-shortly/internal/app/handler"
	"github.com/patraden/ya-practicum-go-shortly/internal/app/logger"
	"github.com/patraden/ya-practicum-go-shortly/internal/app/middleware"
	"github.com/patraden/ya-practicum-go-shortly/internal/app/repository"
	"github.com/patraden/ya-practicum-go-shortly/internal/app/service"
	"github.com/patraden/ya-practicum-go-shortly/internal/app/urlgenerator"
	"github.com/patraden/ya-practicum-go-shortly/internal/app/utils"
)

func setupEndPointTestRouter() http.Handler {
	config := config.DefaultConfig()
	repo := repository.NewInMemoryURLRepository()
	_ = repo.AddURL("shortURL", "https://ya.ru") // Set up necessary test data

	gen := urlgenerator.NewRandURLGenerator(config.URLsize)
	service := service.NewShortenerService(repo, gen, config)
	log := logger.NewLogger(zerolog.InfoLevel).GetLogger()

	return h.NewRouter(service, config, log)
}

func TestEndpoints(t *testing.T) {
	t.Parallel()

	router := setupEndPointTestRouter()

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
	service := service.NewShortenerService(repo, gen, config)
	log := logger.NewLogger(zerolog.InfoLevel).GetLogger()
	handler := h.NewHandler(service, config, log)

	return handler.HandlePost
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
		wantIsURL   bool
	}{
		{"valid URL", "/", `https://ya.ru`, h.ContentTypeText, http.StatusCreated, true},
		{"empty body", "/", ``, h.ContentTypeText, http.StatusBadRequest, false},
		{"invalid URL", "/", `//ya.ru`, h.ContentTypeText, http.StatusBadRequest, false},
		{"invalid path", "/a/b/c", `https://ya.ru`, h.ContentTypeText, http.StatusBadRequest, false},
		{"double slash path", "//", `https://ya.ru`, h.ContentTypeText, http.StatusBadRequest, false},
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

			shortURL, err := io.ReadAll(result.Body)
			require.NoError(t, err)

			assert.Equal(t, test.wantStatus, result.StatusCode)
			assert.Equal(t, test.wantIsURL, utils.IsURL(string(shortURL)))
		})
	}
}

func setupHandlerJSON() http.HandlerFunc {
	config := config.DefaultConfig()
	repo := repository.NewInMemoryURLRepository()
	gen := urlgenerator.NewRandURLGenerator(config.URLsize)
	service := service.NewShortenerService(repo, gen, config)
	log := logger.NewLogger(zerolog.InfoLevel).GetLogger()

	return h.NewHandler(service, config, log).HandlePostJSON
}

func TestHandlePostJSON(t *testing.T) {
	t.Parallel()

	handlePostJSON := setupHandlerJSON()

	tests := []struct {
		name        string
		body        string
		contentType string
		wantStatus  int
		wantIsURL   bool
	}{
		{"valid JSON URL", `{"url":"https://practicum.yandex.ru"}`, h.ContentTypeJSON, http.StatusCreated, true},
		{"malformed JSON", `{"url:"https://practicum.yandex.ru"}`, h.ContentTypeJSON, http.StatusBadRequest, false},
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

			if test.wantIsURL {
				var responseBody map[string]string
				err := json.NewDecoder(result.Body).Decode(&responseBody)
				require.NoError(t, err)
				assert.True(t, utils.IsURL(responseBody["result"]))
			}
		})
	}
}

func TestHandleGet(t *testing.T) {
	t.Parallel()

	config := config.DefaultConfig()
	repo := repository.NewInMemoryURLRepository()
	gen := urlgenerator.NewRandURLGenerator(config.URLsize)
	service := service.NewShortenerService(repo, gen, config)
	longURL := `https://ya.ru`
	serverAddr := `http://localhost:8080/`
	shortURL, _ := service.ShortenURL(longURL)
	log := logger.NewLogger(zerolog.InfoLevel).GetLogger()
	router := h.NewRouter(service, config, log)

	tests := []struct {
		name         string
		path         string
		wantStatus   int
		wantLocation string
	}{
		{"valid short URL", serverAddr + shortURL, http.StatusTemporaryRedirect, longURL},
		{"invalid short URL", serverAddr + "qwerty", http.StatusBadRequest, ""},
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

func setupPostCompressionHandler() http.Handler {
	config := config.DefaultConfig()
	repo := repository.NewInMemoryURLRepository()
	gen := urlgenerator.NewRandURLGenerator(config.URLsize)
	service := service.NewShortenerService(repo, gen, config)
	log := logger.NewLogger(zerolog.InfoLevel).GetLogger()
	handler := http.HandlerFunc(h.NewHandler(service, config, log).HandlePost)

	return middleware.Decompress()(middleware.Compress()(handler))
}

func TestHandlePostCompression(t *testing.T) {
	t.Parallel()

	handler := setupPostCompressionHandler()

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
		for range 10 {
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
				assert.Equal(t, test.isURL, utils.IsURL(string(shortURL)))
			})
		}
	}
}
