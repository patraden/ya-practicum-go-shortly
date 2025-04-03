package grpc_test

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	pb "github.com/patraden/ya-practicum-go-shortly/api/shortener/v1"
	"github.com/patraden/ya-practicum-go-shortly/internal/app/config"
	"github.com/patraden/ya-practicum-go-shortly/internal/app/domain"
	"github.com/patraden/ya-practicum-go-shortly/internal/app/handler"
	"github.com/patraden/ya-practicum-go-shortly/internal/app/logger"
	"github.com/patraden/ya-practicum-go-shortly/internal/app/mock"
	server "github.com/patraden/ya-practicum-go-shortly/internal/app/server/grpc"
)

func setupGRPCServer(t *testing.T) (*gomock.Controller, *mock.MockURLShortener, *server.Server) {
	t.Helper()

	ctrl := gomock.NewController(t)
	mockSrv := mock.NewMockURLShortener(ctrl)
	cfg := &config.Config{BaseURL: "http://base.url", ServerGRPCAddr: "127.0.0.1:50051"}
	log := logger.NewLogger(zerolog.InfoLevel).GetLogger()
	handler, err := handler.NewGRPCURLShortenerHandler(mockSrv, cfg, log)
	require.NoError(t, err)

	srv := server.NewServer(cfg, handler, log)

	return ctrl, mockSrv, srv
}

func TestGRPCServerRunAndShutdown(t *testing.T) {
	t.Parallel()

	ctrl, mockSrv, srv := setupGRPCServer(t)
	defer ctrl.Finish()

	ctx := context.Background()
	wgr := sync.WaitGroup{}

	wgr.Add(1)

	go func() {
		defer wgr.Done()

		err := srv.Run()
		assert.NoError(t, err)
	}()

	time.Sleep(100 * time.Millisecond)

	conn, err := grpc.NewClient("127.0.0.1:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	require.NoError(t, err)
	defer conn.Close()

	client := pb.NewURLShortenerServiceClient(conn)

	mockSrv.EXPECT().
		ShortenURL(gomock.Any(), domain.OriginalURL("https://example.com")).
		Return(domain.Slug("slug1"), nil).
		AnyTimes()

	resp, err := client.ShortenURL(ctx, &pb.ShortenURLRequest{Url: "https://example.com"})
	require.NoError(t, err)
	assert.NotNil(t, resp)

	ctx, cancel := context.WithTimeout(ctx, time.Second)
	defer cancel()

	err = srv.Shutdown(ctx)
	require.NoError(t, err)

	wgr.Wait()
}
