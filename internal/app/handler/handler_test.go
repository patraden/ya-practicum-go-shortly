package handler_test

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/patraden/ya-practicum-go-shortly/internal/app/config"
	h "github.com/patraden/ya-practicum-go-shortly/internal/app/handler"
	"github.com/patraden/ya-practicum-go-shortly/internal/app/middleware"
	"github.com/patraden/ya-practicum-go-shortly/internal/app/repository"
	"github.com/patraden/ya-practicum-go-shortly/internal/app/service"
	"github.com/patraden/ya-practicum-go-shortly/internal/app/urlgenerator"
	"github.com/patraden/ya-practicum-go-shortly/internal/app/utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHandle(t *testing.T) {
	t.Parallel()

	config := config.DefaultConfig()
	repo := repository.NewInMemoryURLRepository()
	err := repo.AddURL("shortURL", "https://ya.ru")
	require.NoError(t, err)

	gen := urlgenerator.NewRandURLGenerator(config.URLsize)
	service := service.NewShortenerService(repo, gen, config)
	router := h.NewRouter(service, config)

	tests := []struct {
		name   string
		method string
		path   string
		body   io.Reader
		want   int
	}{
		{
			name:   "test 1",
			method: http.MethodDelete,
			path:   "/",
			body:   nil,
			want:   http.StatusMethodNotAllowed,
		},
		{
			name:   "test 2",
			method: http.MethodPatch,
			path:   "/",
			body:   nil,
			want:   http.StatusMethodNotAllowed,
		},
		{
			name:   "test 3",
			method: http.MethodPost,
			path:   "/",
			body:   strings.NewReader(`https://ya.ru`),
			want:   http.StatusCreated,
		},
		{
			name:   "test 4",
			method: http.MethodGet,
			path:   "/shortURL",
			body:   nil,
			want:   http.StatusTemporaryRedirect,
		},
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

			if result.StatusCode == http.StatusTemporaryRedirect {
				location := result.Header.Get("Location")
				assert.Equal(t, "https://ya.ru", location)
			}
		})
	}
}

func TestHandlePost(t *testing.T) {
	t.Parallel()

	config := config.DefaultConfig()
	repo := repository.NewInMemoryURLRepository()
	gen := urlgenerator.NewRandURLGenerator(config.URLsize)
	service := service.NewShortenerService(repo, gen, config)
	handler := h.NewHandler(service, config)

	type want struct {
		status int
		isURL  bool
	}
	tests := []struct {
		name        string
		path        string
		body        string
		contentType string
		want        want
	}{
		{
			name:        "test 1",
			body:        `https://ya.ru`,
			path:        `/`,
			contentType: h.ContentTypeText,
			want: want{
				status: http.StatusCreated,
				isURL:  true,
			},
		},
		{
			name:        "test 2",
			body:        ``,
			path:        `/`,
			contentType: h.ContentTypeText,
			want: want{
				status: http.StatusBadRequest,
				isURL:  false,
			},
		},
		{
			name:        "test 3",
			body:        `//ya.ru`,
			path:        `/`,
			contentType: h.ContentTypeText,
			want: want{
				status: http.StatusBadRequest,
				isURL:  false,
			},
		},
		{
			name:        "test 4",
			body:        `https://ya.ru`,
			path:        `/a/b/c`,
			contentType: h.ContentTypeText,
			want: want{
				status: http.StatusBadRequest,
				isURL:  false,
			},
		},
		{
			name:        "test 5",
			body:        `https://ya.ru`,
			path:        `//`,
			contentType: h.ContentTypeText,
			want: want{
				status: http.StatusBadRequest,
				isURL:  false,
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			request := httptest.NewRequest(http.MethodPost, test.path, strings.NewReader(test.body))
			request.Header.Add(h.ContentType, test.contentType)

			w := httptest.NewRecorder()
			h := http.HandlerFunc(handler.HandlePost)
			h(w, request)

			result := w.Result()
			shortURL, err := io.ReadAll(result.Body)

			require.NoError(t, err)
			err = result.Body.Close()
			require.NoError(t, err)

			assert.Equal(t, test.want.status, result.StatusCode)
			assert.Equal(t, test.want.isURL, utils.IsURL(string(shortURL)))
		})
	}
}

func TestHandlePostJSON(t *testing.T) {
	t.Parallel()

	config := config.DefaultConfig()
	repo := repository.NewInMemoryURLRepository()
	gen := urlgenerator.NewRandURLGenerator(config.URLsize)
	service := service.NewShortenerService(repo, gen, config)
	handler := h.NewHandler(service, config).HandlePostJSON

	type want struct {
		status int
		isURL  bool
	}

	tests := []struct {
		name        string
		body        string
		contentType string
		want        want
	}{
		{
			name:        `test 1`,
			body:        `{"url":"https://practicum.yandex.ru"}`,
			contentType: h.ContentTypeJSON,
			want: want{
				status: http.StatusCreated,
				isURL:  true,
			},
		},
		{
			name:        `test 2`,
			body:        `{"url:"https://practicum.yandex.ru"}`,
			contentType: h.ContentTypeJSON,
			want: want{
				status: http.StatusBadRequest,
				isURL:  false,
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			request := httptest.NewRequest(http.MethodPost, `/api/shorten`, strings.NewReader(test.body))
			request.Header.Add(h.ContentType, test.contentType)

			w := httptest.NewRecorder()
			h := http.HandlerFunc(handler)
			h(w, request)

			result := w.Result()

			assert.Equal(t, test.want.status, result.StatusCode)

			if test.want.isURL {
				var responseBody map[string]string

				err := json.NewDecoder(result.Body).Decode(&responseBody)
				require.NoError(t, err)

				err = result.Body.Close()
				require.NoError(t, err)

				assert.Equal(t, test.want.isURL, utils.IsURL(responseBody["result"]))
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
	router := h.NewRouter(service, config)

	type want struct {
		status   int
		location string
	}
	tests := []struct {
		name string
		path string
		want want
	}{
		{
			name: "test 1",
			path: serverAddr + shortURL,
			want: want{
				status:   http.StatusTemporaryRedirect,
				location: longURL,
			},
		},
		{
			name: "test 2",
			path: serverAddr + "qwerty",
			want: want{
				status:   http.StatusBadRequest,
				location: "",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			request := httptest.NewRequest(http.MethodGet, tt.path, nil)
			w := httptest.NewRecorder()

			router.ServeHTTP(w, request)

			result := w.Result()
			defer result.Body.Close()

			assert.Equal(t, tt.want.status, result.StatusCode)

			if result.StatusCode == http.StatusTemporaryRedirect {
				assert.Equal(t, tt.want.location, result.Header.Get("Location"))
			}
		})
	}
}

func TestHandlePostCompression(t *testing.T) {
	t.Parallel()

	config := config.DefaultConfig()
	repo := repository.NewInMemoryURLRepository()
	gen := urlgenerator.NewRandURLGenerator(config.URLsize)
	service := service.NewShortenerService(repo, gen, config)
	h := http.HandlerFunc(h.NewHandler(service, config).HandlePost)
	hd := middleware.Compress()(h)
	hdc := middleware.Decompress()(hd)

	type want struct {
		status int
		isURL  bool
	}
	tests := []struct {
		name            string
		contentEncoding string
		acceptEncoding  string
		want            want
	}{
		{
			name:            "test 1",
			contentEncoding: "gzip",
			acceptEncoding:  "",
			want: want{
				status: http.StatusCreated,
				isURL:  true,
			},
		},
		{
			name:            "test 2",
			contentEncoding: "deflate",
			acceptEncoding:  "",
			want: want{
				status: http.StatusCreated,
				isURL:  true,
			},
		},
		{
			name:            "test 3",
			contentEncoding: "gzip",
			acceptEncoding:  "deflate",
			want: want{
				status: http.StatusCreated,
				isURL:  true,
			},
		},
		{
			name:            "test 4",
			contentEncoding: "deflate",
			acceptEncoding:  "gzip",
			want: want{
				status: http.StatusCreated,
				isURL:  true,
			},
		},
		{
			name:            "test 5",
			contentEncoding: "",
			acceptEncoding:  "",
			want: want{
				status: http.StatusCreated,
				isURL:  true,
			},
		},
	}

	for _, test := range tests {
		for range 10 {
			t.Run(test.name, func(t *testing.T) {
				t.Parallel()

				url := utils.RandURL()

				data, err := utils.Compress([]byte(url), test.contentEncoding)
				require.NoError(t, err)

				w := httptest.NewRecorder()
				r := httptest.NewRequest(http.MethodPost, "/", bytes.NewReader(data))

				r.Header.Set("Content-Encoding", test.contentEncoding)
				r.Header.Set("Accept-Encoding", test.acceptEncoding)

				hdc.ServeHTTP(w, r)

				result := w.Result()
				compressedURL, err := io.ReadAll(result.Body)
				require.NoError(t, err)

				shortURL, err := utils.Decompress(compressedURL, test.acceptEncoding)
				require.NoError(t, err)

				err = result.Body.Close()
				require.NoError(t, err)

				assert.Equal(t, test.want.status, result.StatusCode)
				assert.Equal(t, test.want.isURL, utils.IsURL(string(shortURL)))
			})
		}
	}
}
