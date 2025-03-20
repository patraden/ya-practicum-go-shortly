package handler_test

import (
	"bytes"
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/mailru/easyjson"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	"github.com/patraden/ya-practicum-go-shortly/internal/app/config"
	"github.com/patraden/ya-practicum-go-shortly/internal/app/dto"
	"github.com/patraden/ya-practicum-go-shortly/internal/app/handler"
	"github.com/patraden/ya-practicum-go-shortly/internal/app/logger"
	"github.com/patraden/ya-practicum-go-shortly/internal/app/mock"
)

func setupDeleteHandler(t *testing.T) (*gomock.Controller, *mock.MockURLRemover, *handler.DeleteHandler) {
	t.Helper()

	ctrl := gomock.NewController(t)
	mockSrv := mock.NewMockURLRemover(ctrl)
	log := logger.NewLogger(zerolog.InfoLevel).GetLogger()
	config := &config.Config{BaseURL: "http://base.url"}
	h := handler.NewDeleteHandler(mockSrv, config, log)

	return ctrl, mockSrv, h
}

func TestHandleDelUserURLsSuccess(t *testing.T) {
	t.Parallel()

	ctrl, mockRemover, hlr := setupDeleteHandler(t)
	defer ctrl.Finish()

	slugs := dto.UserSlugBatch{"slug1", "slug2"}
	jsonBody, err := easyjson.Marshal(slugs)
	require.NoError(t, err)

	req := httptest.NewRequest(http.MethodDelete, "/api/user/urls", bytes.NewReader(jsonBody))
	req = req.WithContext(context.Background())

	rec := httptest.NewRecorder()

	mockRemover.EXPECT().
		RemoveUserSlugs(gomock.Any(), slugs).
		Return(nil)

	hlr.HandleDelUserURLs(rec, req)

	res := rec.Result()
	defer res.Body.Close()

	assert.Equal(t, http.StatusAccepted, res.StatusCode)
	assert.Equal(t, handler.ContentTypeText, res.Header.Get("Content-Type"))
}

func TestHandleDelUserURLsBadRequest(t *testing.T) {
	t.Parallel()

	_, _, hlr := setupDeleteHandler(t)

	req := httptest.NewRequest(http.MethodDelete, "/api/user/urls", bytes.NewReader([]byte("{invalid json}")))
	req = req.WithContext(context.Background())

	rec := httptest.NewRecorder()

	hlr.HandleDelUserURLs(rec, req)

	res := rec.Result()
	defer res.Body.Close()

	assert.Equal(t, http.StatusBadRequest, res.StatusCode)
}

func TestHandleDelUserURLsInternalError(t *testing.T) {
	t.Parallel()

	ctrl, mockRemover, hlr := setupDeleteHandler(t)
	defer ctrl.Finish()

	slugs := dto.UserSlugBatch{"slug1", "slug2"}
	jsonBody, err := easyjson.Marshal(slugs)
	require.NoError(t, err)

	req := httptest.NewRequest(http.MethodDelete, "/api/user/urls", bytes.NewReader(jsonBody))
	req = req.WithContext(context.Background())

	rec := httptest.NewRecorder()

	mockRemover.EXPECT().
		RemoveUserSlugs(gomock.Any(), slugs).
		Return(assert.AnError)

	hlr.HandleDelUserURLs(rec, req)

	res := rec.Result()
	defer res.Body.Close()

	assert.Equal(t, http.StatusInternalServerError, res.StatusCode)
}
