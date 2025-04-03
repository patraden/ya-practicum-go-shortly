package handler_test

import (
	"context"
	"testing"

	"github.com/rs/zerolog"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	pb "github.com/patraden/ya-practicum-go-shortly/api/shortener/v1"
	"github.com/patraden/ya-practicum-go-shortly/internal/app/config"
	"github.com/patraden/ya-practicum-go-shortly/internal/app/domain"
	e "github.com/patraden/ya-practicum-go-shortly/internal/app/domain/errors"
	"github.com/patraden/ya-practicum-go-shortly/internal/app/handler"
	"github.com/patraden/ya-practicum-go-shortly/internal/app/logger"
	"github.com/patraden/ya-practicum-go-shortly/internal/app/mock"
)

func setupGRPCShortenerHandler(t *testing.T) (
	*gomock.Controller,
	*mock.MockURLShortener,
	*handler.GRPCShortenerHandler,
) {
	t.Helper()

	ctrl := gomock.NewController(t)
	mockSrv := mock.NewMockURLShortener(ctrl)
	log := logger.NewLogger(zerolog.InfoLevel).GetLogger()
	config := &config.Config{BaseURL: "http://base.url"}
	h, err := handler.NewGRPCURLShortenerHandler(mockSrv, config, log)
	require.NoError(t, err)

	return ctrl, mockSrv, h
}

func TestGRPCSShortenURL(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		url         string
		mockReturn  domain.Slug
		mockError   error
		expectedErr codes.Code
		expectedURL string
	}{
		{"Success", "https://example.com", domain.Slug("abcd1234"), nil, codes.OK, "http://base.url/abcd1234"},
		{"Invalid URL format", "invalid-url", "", nil, codes.InvalidArgument, ""},
		{"Already Exists", "https://duplicate.com", "", e.ErrOriginalExists, codes.AlreadyExists, ""},
		{"Internal Error", "https://example.com", "", e.ErrTestGeneral, codes.Internal, ""},
	}

	for _, ttc := range tests {
		t.Run(ttc.name, func(t *testing.T) {
			t.Parallel()

			ctrl, mockSrv, h := setupGRPCShortenerHandler(t)
			defer ctrl.Finish()

			if ttc.mockError != nil || ttc.mockReturn != "" {
				mockSrv.EXPECT().
					ShortenURL(gomock.Any(), domain.OriginalURL(ttc.url)).
					Return(ttc.mockReturn, ttc.mockError).
					AnyTimes()
			}

			resp, err := h.ShortenURL(context.Background(), &pb.ShortenURLRequest{Url: ttc.url})

			if ttc.expectedErr == codes.OK {
				require.NoError(t, err)
				require.NotNil(t, resp)
				require.Equal(t, ttc.expectedURL, resp.GetSlug())
			} else {
				require.Error(t, err)
				require.Equal(t, ttc.expectedErr, status.Code(err))
			}
		})
	}
}

func TestGetOriginalURL(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		slug        string
		mockReturn  domain.OriginalURL
		mockError   error
		expectedErr codes.Code
		expectedURL string
	}{
		{"Success", "abcd1234", domain.OriginalURL("https://example.com"), nil, codes.OK, "https://example.com"},
		{"Slug Too Short", "abc", "", nil, codes.InvalidArgument, ""},
		{"Invalid Slug", "invalid-slug", "", e.ErrSlugInvalid, codes.InvalidArgument, ""},
		{"Slug Not Found", "notfound123", "", e.ErrSlugNotFound, codes.NotFound, ""},
		{"Slug Deleted", "deleted123", "", e.ErrSlugDeleted, codes.NotFound, ""},
		{"Internal Error", "error123", "", e.ErrShortenerInternal, codes.Internal, ""},
	}

	for _, ttc := range tests {
		t.Run(ttc.name, func(t *testing.T) {
			t.Parallel()

			ctrl, mockSrv, h := setupGRPCShortenerHandler(t)
			defer ctrl.Finish()

			if ttc.mockError != nil || ttc.mockReturn != "" {
				mockSrv.EXPECT().
					GetOriginalURL(gomock.Any(), domain.Slug(ttc.slug)).
					Return(ttc.mockReturn, ttc.mockError).
					AnyTimes()
			}

			resp, err := h.GetOriginalURL(context.Background(), &pb.GetOriginalURLRequest{Slug: ttc.slug})

			if ttc.expectedErr == codes.OK {
				require.NoError(t, err)
				require.NotNil(t, resp)
				require.Equal(t, ttc.expectedURL, resp.GetUrl())
			} else {
				require.Error(t, err)
				require.Equal(t, ttc.expectedErr, status.Code(err))
			}
		})
	}
}
