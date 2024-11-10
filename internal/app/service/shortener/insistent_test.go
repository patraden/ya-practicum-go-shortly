package shortener_test

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	"github.com/patraden/ya-practicum-go-shortly/internal/app/config"
	"github.com/patraden/ya-practicum-go-shortly/internal/app/domain"
	e "github.com/patraden/ya-practicum-go-shortly/internal/app/errors"
	"github.com/patraden/ya-practicum-go-shortly/internal/app/mock"
	"github.com/patraden/ya-practicum-go-shortly/internal/app/repository"
	"github.com/patraden/ya-practicum-go-shortly/internal/app/service/shortener"
	"github.com/patraden/ya-practicum-go-shortly/internal/app/service/urlgenerator"
	"github.com/patraden/ya-practicum-go-shortly/internal/app/utils"
)

const numStoreRestoreTests = 100

func RunShortenURLCollisionsTests(t *testing.T, config *config.Config, repo repository.URLRepository) {
	t.Helper()

	ctrl := gomock.NewController(t)
	mockURLGen := mock.NewMockURLGenerator(ctrl)
	// 1 call to generate initial shortURL
	calls := int(config.URLGenTimeout/config.URLGenRetryInterval) + 1

	mockURLGen.
		EXPECT().
		GenerateSlug(gomock.Any(), gomock.Any()).
		Return(domain.Slug("shortURL")).
		Times(calls)

	srv := shortener.NewInsistentShortener(repo, mockURLGen, config)

	t.Run("url collisions test", func(t *testing.T) {
		t.Parallel()

		_, err := srv.ShortenURL(context.Background(), domain.OriginalURL("abc"))

		require.NoError(t, err)

		start := time.Now()
		_, err = srv.ShortenURL(context.Background(), domain.OriginalURL("cba"))

		require.ErrorIs(t, err, e.ErrServiceCollision)
		assert.LessOrEqual(t, time.Since(start), config.URLGenTimeout)
	})
}

func RunStoreRestoreTests(t *testing.T, srv shortener.URLShortener) {
	t.Helper()

	t.Run("ShortenURL", func(t *testing.T) {
		t.Parallel()

		original := domain.OriginalURL(utils.RandURL())
		mapURL, err := srv.ShortenURL(context.Background(), original)

		require.NoError(t, err)

		newMapURL, err := srv.GetOriginalURL(context.Background(), mapURL.Slug)
		require.NoError(t, err)
		assert.Equal(t, original, newMapURL.OriginalURL)
	})
}

func RunShortenURLErrTests(t *testing.T, config *config.Config, gen urlgenerator.URLGenerator) {
	t.Helper()

	ctrl := gomock.NewController(t)
	mockURLRepo := mock.NewMockURLRepository(ctrl)

	// repo is unavailable
	gomock.InOrder(
		mockURLRepo.
			EXPECT().
			GetURLMapping(gomock.Any(), gomock.Any()).
			Return(nil, e.ErrTest).
			Times(1),

		mockURLRepo.
			EXPECT().
			GetURLMapping(gomock.Any(), gomock.Any()).
			Return(nil, e.ErrRepoNotFound).
			Times(1),

		mockURLRepo.
			EXPECT().
			AddURLMapping(gomock.Any(), gomock.Any()).
			Return(e.ErrTest).
			Times(1),
	)

	srv := shortener.NewInsistentShortener(mockURLRepo, gen, config)

	t.Run("ShortenURLErr", func(t *testing.T) {
		t.Parallel()

		// first call should have problems to get from repo.
		originalURL := domain.OriginalURL("http://localhost:8080")
		_, err := srv.ShortenURL(context.Background(), originalURL)
		require.ErrorIs(t, err, e.ErrServiceInternal)

		// second call should have problems to add to repo.
		originalURL = domain.OriginalURL("http://localhost:8181")
		_, err = srv.ShortenURL(context.Background(), originalURL)
		require.ErrorIs(t, err, e.ErrServiceInternal)
	})
}

func RunGetOriginalURLErrTests(t *testing.T, config *config.Config, gen urlgenerator.URLGenerator) {
	t.Helper()

	ctrl := gomock.NewController(t)
	mockURLRepo := mock.NewMockURLRepository(ctrl)

	gomock.InOrder(
		mockURLRepo.
			EXPECT().
			GetURLMapping(gomock.Any(), gomock.Any()).
			Return(nil, e.ErrTest).
			Times(1),

		mockURLRepo.
			EXPECT().
			GetURLMapping(gomock.Any(), gomock.Any()).
			Return(nil, e.ErrRepoNotFound).
			Times(1),
	)

	srv := shortener.NewInsistentShortener(mockURLRepo, gen, config)
	badURL := domain.Slug(utils.RandomString(config.URLsize + 1))

	_, err := srv.GetOriginalURL(context.Background(), badURL)
	require.ErrorIs(t, err, e.ErrServiceInvalid)

	_, err = srv.GetOriginalURL(context.Background(), domain.Slug("shortURL"))
	require.ErrorIs(t, err, e.ErrServiceInternal)

	_, err = srv.GetOriginalURL(context.Background(), domain.Slug("shortURL"))
	require.ErrorIs(t, err, e.ErrRepoNotFound)
}

func TestURLShortener(t *testing.T) {
	t.Parallel()

	config := config.DefaultConfig()
	repo := repository.NewInMemoryURLRepository()
	gen := urlgenerator.NewRandURLGenerator(config.URLsize)
	srv := shortener.NewInsistentShortener(repo, gen, config)

	// store and restore good urls.
	for range numStoreRestoreTests {
		RunStoreRestoreTests(t, srv)
	}

	// colisions and retries.
	RunShortenURLCollisionsTests(t, config, repo)

	// internal errors.
	RunShortenURLErrTests(t, config, gen)
	RunGetOriginalURLErrTests(t, config, gen)
}
