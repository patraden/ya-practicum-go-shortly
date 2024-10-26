package handler_test

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/patraden/ya-practicum-go-shortly/internal/app/config"
	h "github.com/patraden/ya-practicum-go-shortly/internal/app/handler"
	"github.com/patraden/ya-practicum-go-shortly/internal/app/mock"
	"github.com/patraden/ya-practicum-go-shortly/internal/app/service"
	"github.com/patraden/ya-practicum-go-shortly/internal/app/utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestHandle(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	mockService := mock.NewMockURLShortener(ctrl)
	mockService.
		EXPECT().
		ShortenURL(gomock.Eq("https://ya.ru")).
		Return("shortURL", nil).
		AnyTimes()

	mockService.
		EXPECT().
		GetOriginalURL(gomock.Eq("shortURL")).
		Return("https://ya.ru", nil).
		AnyTimes()

	config := config.DefaultConfig()
	router := h.NewRouter(mockService, config)

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
	service := service.NewInMemoryShortenerService(config)
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

func TestHandleGet(t *testing.T) {
	t.Parallel()

	config := config.DefaultConfig()
	service := service.NewInMemoryShortenerService(config)
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
