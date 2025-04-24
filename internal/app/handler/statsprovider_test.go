package handler_test

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	"github.com/patraden/ya-practicum-go-shortly/internal/app/config"
	e "github.com/patraden/ya-practicum-go-shortly/internal/app/domain/errors"
	"github.com/patraden/ya-practicum-go-shortly/internal/app/dto"
	"github.com/patraden/ya-practicum-go-shortly/internal/app/handler"
	"github.com/patraden/ya-practicum-go-shortly/internal/app/logger"
	"github.com/patraden/ya-practicum-go-shortly/internal/app/mock"
)

type failingWriter struct {
	http.ResponseWriter
	StatusCode int
}

func (w *failingWriter) Write(_ []byte) (int, error) {
	return 0, e.ErrTestGeneral
}

func (w *failingWriter) WriteHeader(statusCode int) {
	w.StatusCode = statusCode
}

func (w *failingWriter) Header() http.Header {
	return http.Header{}
}

func setupStatsProviderHandler(t *testing.T) (
	*gomock.Controller,
	*mock.MockStatsProvider,
	*handler.StatsProviderHandler,
) {
	t.Helper()

	ctrl := gomock.NewController(t)
	mockSrv := mock.NewMockStatsProvider(ctrl)
	log := logger.NewLogger(zerolog.InfoLevel).GetLogger()
	config := config.DefaultConfig()
	h := handler.NewStatsProviderHandler(mockSrv, config, log)

	return ctrl, mockSrv, h
}

func TestHandleGetStats(t *testing.T) {
	t.Parallel()

	ctrl, mockSrv, handler := setupStatsProviderHandler(t)
	defer ctrl.Finish()

	stats := &dto.RepoStats{
		CountSlugs: int64(2),
		CountUsers: int64(1),
	}

	t.Run("successful request", func(t *testing.T) {
		mockSrv.EXPECT().GetStats(gomock.Any()).Return(stats, nil)

		req := httptest.NewRequest(http.MethodGet, "/api/internal/stats", nil)
		w := httptest.NewRecorder()

		handler.HandleGetStats(w, req)
		res := w.Result()

		defer res.Body.Close()

		body, _ := io.ReadAll(res.Body)

		assert.Equal(t, http.StatusOK, res.StatusCode)
		assert.JSONEq(t, `{"urls":2,"users":1}`, string(body))
	})

	t.Run("failed service", func(t *testing.T) {
		mockSrv.EXPECT().GetStats(gomock.Any()).Return(nil, e.ErrTestGeneral)

		req := httptest.NewRequest(http.MethodGet, "/api/internal/stats", nil)
		w := httptest.NewRecorder()

		handler.HandleGetStats(w, req)
		res := w.Result()

		defer res.Body.Close()

		assert.Equal(t, http.StatusInternalServerError, res.StatusCode)
	})

	t.Run("failed response write", func(t *testing.T) {
		mockSrv.EXPECT().GetStats(gomock.Any()).Return(stats, nil)

		req := httptest.NewRequest(http.MethodGet, "/api/internal/stats", nil)
		w := httptest.NewRecorder()
		fw := &failingWriter{
			ResponseWriter: w,
			StatusCode:     http.StatusOK,
		}

		handler.HandleGetStats(fw, req)

		assert.Equal(t, http.StatusInternalServerError, fw.StatusCode)
	})
}
