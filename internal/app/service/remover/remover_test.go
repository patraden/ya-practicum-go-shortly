package remover_test

import (
	"context"
	"strconv"
	"sync/atomic"
	"testing"
	"time"

	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	"github.com/patraden/ya-practicum-go-shortly/internal/app/domain"
	"github.com/patraden/ya-practicum-go-shortly/internal/app/dto"
	"github.com/patraden/ya-practicum-go-shortly/internal/app/logger"
	"github.com/patraden/ya-practicum-go-shortly/internal/app/middleware"
	"github.com/patraden/ya-practicum-go-shortly/internal/app/mock"
	"github.com/patraden/ya-practicum-go-shortly/internal/app/service/remover"
)

func TestAsyncRemover(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	mockRepo := mock.NewMockURLRepository(ctrl)
	log := logger.NewLogger(zerolog.DebugLevel).GetLogger()

	var taskCounter int32

	mockRepo.EXPECT().
		DelUserURLMappings(gomock.Any(), gomock.Any()).
		DoAndReturn(func(_ context.Context, tasks *[]dto.UserSlug) error {
			for _, task := range *tasks {
				time.Sleep(time.Millisecond)
				atomic.AddInt32(&taskCounter, 1)
				log.Info().
					Str("Slug", task.Slug.String()).
					Str("UserID", task.UserID.String()).
					Msg("Deleted slug")
			}

			return nil
		}).
		AnyTimes()

	user := domain.NewUserID()
	slugs := []domain.Slug{}

	expectedTasks := 100
	for i := range expectedTasks {
		slugs = append(slugs, domain.Slug("slug"+strconv.Itoa(i)))
	}

	ctx := context.WithValue(context.Background(), middleware.UserIDKey, user)
	remover := remover.NewAsyncRemover(100*time.Millisecond, mockRepo, log)
	remover.Start()
	assert.True(t, remover.IsRunning())

	err := remover.RemoveUserSlugs(ctx, slugs)
	require.NoError(t, err)

	remover.Stop(ctx)
	assert.Equal(t, expectedTasks, int(taskCounter))
}