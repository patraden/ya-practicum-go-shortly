package handler_test

import (
	"bytes"
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/mailru/easyjson"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	"github.com/patraden/ya-practicum-go-shortly/internal/app/config"
	"github.com/patraden/ya-practicum-go-shortly/internal/app/domain"
	e "github.com/patraden/ya-practicum-go-shortly/internal/app/domain/errors"
	"github.com/patraden/ya-practicum-go-shortly/internal/app/dto"
	"github.com/patraden/ya-practicum-go-shortly/internal/app/handler"
	"github.com/patraden/ya-practicum-go-shortly/internal/app/logger"
	"github.com/patraden/ya-practicum-go-shortly/internal/app/mock"
)

func setupHandler(t *testing.T) (*gomock.Controller, *mock.MockURLShortener, *handler.ShortenerHandler) {
	t.Helper()

	ctrl := gomock.NewController(t)
	mockSrv := mock.NewMockURLShortener(ctrl)
	log := logger.NewLogger(zerolog.InfoLevel).GetLogger()
	config := &config.Config{BaseURL: "http://base.url"}
	h := handler.NewShortenerHandler(mockSrv, config, log)

	return ctrl, mockSrv, h
}

type testCaseGetOriginalURL struct {
	name         string
	shortURL     string
	mockBehavior func()
	expectedCode int
	expectedBody string
}

func TestHandleGetOriginalURL(t *testing.T) {
	t.Parallel()

	ctrl, mockSrv, handler := setupHandler(t)
	defer ctrl.Finish()

	tests := []testCaseGetOriginalURL{
		{
			name:     "Successful Redirect",
			shortURL: "shortURL",
			mockBehavior: func() {
				mockSrv.EXPECT().GetOriginalURL(gomock.Any(), domain.Slug("shortURL")).
					Return(domain.OriginalURL("https://ya.ru"), nil).Times(1)
			},
			expectedCode: http.StatusTemporaryRedirect,
			expectedBody: "",
		},
		{
			name:     "Slug Not Found",
			shortURL: "shortURL",
			mockBehavior: func() {
				mockSrv.EXPECT().GetOriginalURL(gomock.Any(), domain.Slug("shortURL")).
					Return(domain.OriginalURL(""), e.ErrSlugNotFound)
			},
			expectedCode: http.StatusNotFound,
			expectedBody: "slug not found",
		},
		{
			name:     "Slug Invalid",
			shortURL: "shortURL",
			mockBehavior: func() {
				mockSrv.EXPECT().GetOriginalURL(gomock.Any(), domain.Slug("shortURL")).
					Return(domain.OriginalURL(""), e.ErrSlugInvalid)
			},
			expectedCode: http.StatusBadRequest,
			expectedBody: "invalid slug",
		},
		{
			name:     "Internal Error",
			shortURL: "shortURL",
			mockBehavior: func() {
				mockSrv.EXPECT().GetOriginalURL(gomock.Any(), domain.Slug("shortURL")).
					Return(domain.OriginalURL(""), e.ErrShortenerInternal)
			},
			expectedCode: http.StatusInternalServerError,
			expectedBody: "internal error",
		},
	}

	for _, test := range tests {
		runHandleGetOriginalURL(t, handler, test)
	}
}

func runHandleGetOriginalURL(t *testing.T, handler *handler.ShortenerHandler, test testCaseGetOriginalURL) {
	t.Helper()

	t.Run(test.name, func(t *testing.T) {
		test.mockBehavior()

		req := httptest.NewRequest(http.MethodGet, "/{shortURL}/", nil)
		rctx := chi.NewRouteContext()
		rctx.URLParams.Add("shortURL", test.shortURL)
		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

		w := httptest.NewRecorder()

		handler.HandleGetOriginalURL(w, req)

		res := w.Result()
		defer res.Body.Close()

		assert.Equal(t, test.expectedCode, res.StatusCode)

		if test.expectedBody != "" {
			body, _ := io.ReadAll(res.Body)
			assert.Contains(t, string(body), test.expectedBody)
		}
	})
}

func TestHandleShortenURL(t *testing.T) {
	t.Parallel()

	ctrl, mockSrv, handler := setupHandler(t)
	defer ctrl.Finish()

	tests := []struct {
		name         string
		body         string
		mockBehavior func()
		expectedCode int
		expectedBody string
	}{
		{
			name: "Successful Shorten URL",
			body: "https://example.com",
			mockBehavior: func() {
				mockSrv.EXPECT().ShortenURL(gomock.Any(), domain.OriginalURL("https://example.com")).
					Return(domain.Slug("shortURL"), nil).Times(1)
			},
			expectedCode: http.StatusCreated,
			expectedBody: "http://base.url/shortURL",
		},
		{
			name: "Original Exists Conflict",
			body: "https://example.com",
			mockBehavior: func() {
				mockSrv.EXPECT().ShortenURL(gomock.Any(), domain.OriginalURL("https://example.com")).
					Return(domain.Slug("shortURL"), e.ErrOriginalExists)
			},
			expectedCode: http.StatusConflict,
			expectedBody: "http://base.url/shortURL",
		},
		{
			name:         "Invalid URL",
			body:         "invalid-url",
			mockBehavior: func() {},
			expectedCode: http.StatusBadRequest,
			expectedBody: "bad request",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockBehavior()

			req := httptest.NewRequest(http.MethodPost, "/", io.NopCloser(strings.NewReader(tt.body)))
			w := httptest.NewRecorder()

			handler.HandleShortenURL(w, req)

			res := w.Result()
			defer res.Body.Close()

			assert.Equal(t, tt.expectedCode, res.StatusCode)

			body, _ := io.ReadAll(res.Body)
			assert.Contains(t, string(body), tt.expectedBody)
		})
	}
}

type testCaseShortenURLJSON struct {
	name         string
	body         string
	mockBehavior func()
	expectedCode int
	expectedBody string
}

func TestHandleShortenURLJSON(t *testing.T) {
	t.Parallel()

	ctrl, mockSrv, handler := setupHandler(t)
	defer ctrl.Finish()

	tests := []testCaseShortenURLJSON{
		{
			name: "Successful Shorten URL JSON",
			body: `{"url": "https://example.com"}`,
			mockBehavior: func() {
				mockSrv.EXPECT().ShortenURL(gomock.Any(), domain.OriginalURL("https://example.com")).
					Return(domain.Slug("shortURL"), nil).Times(1)
			},
			expectedCode: http.StatusCreated,
			expectedBody: `"result":"http://base.url/shortURL"`,
		},
		{
			name: "Original Exists Conflict JSON",
			body: `{"url": "https://example.com"}`,
			mockBehavior: func() {
				mockSrv.EXPECT().ShortenURL(gomock.Any(), domain.OriginalURL("https://example.com")).
					Return(domain.Slug("shortURL"), e.ErrOriginalExists)
			},
			expectedCode: http.StatusConflict,
			expectedBody: `"result":"http://base.url/shortURL"`,
		},
		{
			name:         "Invalid JSON",
			body:         `invalid json`,
			mockBehavior: func() {},
			expectedCode: http.StatusBadRequest,
			expectedBody: "invalid json",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			test.mockBehavior()

			req := httptest.NewRequest(http.MethodPost, "/", io.NopCloser(strings.NewReader(test.body)))
			req.Header.Set("Content-Type", "application/json")

			w := httptest.NewRecorder()

			handler.HandleShortenURLJSON(w, req)

			res := w.Result()
			defer res.Body.Close()

			assert.Equal(t, test.expectedCode, res.StatusCode)

			body, _ := io.ReadAll(res.Body)
			assert.Contains(t, string(body), test.expectedBody)
		})
	}
}

type testCaseBatchShortenURLJSON struct {
	name         string
	inputBatch   dto.OriginalURLBatch
	mockBehavior func()
	expectedCode int
	expectedBody dto.SlugBatch
}

func TestHandleBatchShortenURLJSON(t *testing.T) {
	t.Parallel()

	ctrl, mockSrv, handler := setupHandler(t)
	defer ctrl.Finish()

	tests := []testCaseBatchShortenURLJSON{
		{
			name: "Successful Batch Shorten",
			inputBatch: dto.OriginalURLBatch{
				{CorrelationID: "1", OriginalURL: domain.OriginalURL("https://example1.com")},
				{CorrelationID: "2", OriginalURL: domain.OriginalURL("https://example2.com")},
			},
			mockBehavior: func() {
				mockSrv.EXPECT().ShortenURLBatch(gomock.Any(), &dto.OriginalURLBatch{
					{CorrelationID: "1", OriginalURL: domain.OriginalURL("https://example1.com")},
					{CorrelationID: "2", OriginalURL: domain.OriginalURL("https://example2.com")},
				}).Return(&dto.SlugBatch{
					{CorrelationID: "1", Slug: domain.Slug("short1")},
					{CorrelationID: "2", Slug: domain.Slug("short2")},
				}, nil).Times(1)
			},
			expectedCode: http.StatusCreated,
			expectedBody: dto.SlugBatch{
				{CorrelationID: "1", Slug: domain.Slug("http://base.url/short1")},
				{CorrelationID: "2", Slug: domain.Slug("http://base.url/short2")},
			},
		},
		{
			name: "Internal Error",
			inputBatch: dto.OriginalURLBatch{
				{CorrelationID: "1", OriginalURL: domain.OriginalURL("https://example1.com")},
			},
			mockBehavior: func() {
				mockSrv.EXPECT().ShortenURLBatch(gomock.Any(), gomock.Any()).
					Return(nil, e.ErrShortenerInternal)
			},
			expectedCode: http.StatusInternalServerError,
			expectedBody: nil,
		},
	}

	for _, test := range tests {
		runBatchShortenURLJSON(t, handler, test)
	}
}

func runBatchShortenURLJSON(t *testing.T, handler *handler.ShortenerHandler, test testCaseBatchShortenURLJSON) {
	t.Helper()

	t.Run(test.name, func(t *testing.T) {
		test.mockBehavior()

		reqBody, _ := easyjson.Marshal(test.inputBatch)
		req := httptest.NewRequest(http.MethodPost, "/api/shorten/batch", io.NopCloser(bytes.NewReader(reqBody)))
		w := httptest.NewRecorder()

		handler.HandleBatchShortenURLJSON(w, req)

		res := w.Result()
		defer res.Body.Close()

		assert.Equal(t, test.expectedCode, res.StatusCode)

		if test.expectedBody != nil {
			var response dto.SlugBatch
			body, _ := io.ReadAll(res.Body)
			_ = easyjson.Unmarshal(body, &response)

			for i, expectedSlug := range test.expectedBody {
				assert.Equal(t, expectedSlug.CorrelationID, response[i].CorrelationID)
				assert.Equal(t, expectedSlug.Slug, response[i].Slug)
			}
		}
	})
}

type testCaseGetUserURLs struct {
	name         string
	mockBehavior func()
	expectedCode int
	expectedBody string
}

func TestHandleGetUserURLs(t *testing.T) {
	t.Parallel()

	ctrl, mockSrv, handler := setupHandler(t)
	defer ctrl.Finish()

	tests := []testCaseGetUserURLs{
		{
			name: "Successful URL Retrieval",
			mockBehavior: func() {
				mockSrv.EXPECT().GetUserURLs(gomock.Any()).
					Return(&dto.URLPairBatch{
						{Slug: "http://base.url/short1", OriginalURL: "https://example1.com"},
						{Slug: "http://base.url/short2", OriginalURL: "https://example2.com"},
					}, nil)
			},
			expectedCode: http.StatusOK,
			expectedBody: `[{"short_url":"http://base.url/short1","original_url":"https://example1.com"},
											{"short_url":"http://base.url/short2","original_url":"https://example2.com"}]`,
		},
		{
			name: "No URLs Found",
			mockBehavior: func() {
				mockSrv.EXPECT().GetUserURLs(gomock.Any()).Return(&dto.URLPairBatch{}, e.ErrUserNotFound)
			},
			expectedCode: http.StatusNoContent,
			expectedBody: ``,
		},
		{
			name: "Internal Server Error",
			mockBehavior: func() {
				mockSrv.EXPECT().GetUserURLs(gomock.Any()).Return(&dto.URLPairBatch{}, e.ErrShortenerInternal)
			},
			expectedCode: http.StatusInternalServerError,
			expectedBody: ``,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			test.mockBehavior()

			req := httptest.NewRequest(http.MethodGet, "/api/user/urls", nil)
			rctx := chi.NewRouteContext()
			req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

			w := httptest.NewRecorder()

			handler.HandleGetUserURLs(w, req)

			res := w.Result()
			defer res.Body.Close()

			assert.Equal(t, test.expectedCode, res.StatusCode)

			if test.expectedBody != "" {
				body, _ := io.ReadAll(res.Body)
				assert.JSONEq(t, test.expectedBody, string(body))
			}
		})
	}
}
