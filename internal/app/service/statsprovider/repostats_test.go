package statsprovider_test

import (
	"context"
	"testing"

	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	e "github.com/patraden/ya-practicum-go-shortly/internal/app/domain/errors"
	"github.com/patraden/ya-practicum-go-shortly/internal/app/dto"
	"github.com/patraden/ya-practicum-go-shortly/internal/app/logger"
	"github.com/patraden/ya-practicum-go-shortly/internal/app/mock"
	"github.com/patraden/ya-practicum-go-shortly/internal/app/service/statsprovider"
)

func setupRepoStatsProviderTest(t *testing.T) (
	*gomock.Controller,
	*mock.MockURLRepository,
	*statsprovider.RepoStatsProvider,
) {
	t.Helper()
	ctrl := gomock.NewController(t)
	repo := mock.NewMockURLRepository(ctrl)
	log := logger.NewLogger(zerolog.DebugLevel).GetLogger()
	svc := statsprovider.NewRepoStatsProvider(repo, log)

	return ctrl, repo, svc
}

func TestStatsProvider(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	ctrl, repo, svc := setupRepoStatsProviderTest(t)

	defer ctrl.Finish()

	t.Run("Success test", func(t *testing.T) {
		expectedStats := &dto.RepoStats{
			CountSlugs: int64(10),
			CountUsers: int64(5),
		}

		repo.EXPECT().GetStats(gomock.Any()).Return(expectedStats, nil)

		actualStats, err := svc.GetStats(ctx)
		require.NoError(t, err)

		assert.Equal(t, expectedStats.CountSlugs, actualStats.CountSlugs)
		assert.Equal(t, expectedStats.CountUsers, actualStats.CountUsers)
	})

	t.Run("Failure test", func(t *testing.T) {
		repo.EXPECT().GetStats(gomock.Any()).Return(nil, e.ErrTestGeneral)

		actualStats, err := svc.GetStats(ctx)
		require.ErrorIs(t, err, e.ErrStatsProviderInternal)

		require.Nil(t, actualStats)
	})
}
