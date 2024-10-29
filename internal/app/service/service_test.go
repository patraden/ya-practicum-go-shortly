package service_test

import (
	"testing"
	"time"

	"github.com/patraden/ya-practicum-go-shortly/internal/app/config"
	e "github.com/patraden/ya-practicum-go-shortly/internal/app/errors"
	"github.com/patraden/ya-practicum-go-shortly/internal/app/mock"
	"github.com/patraden/ya-practicum-go-shortly/internal/app/repository"
	"github.com/patraden/ya-practicum-go-shortly/internal/app/service"
	"github.com/patraden/ya-practicum-go-shortly/internal/app/urlgenerator"
	"github.com/patraden/ya-practicum-go-shortly/internal/app/utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
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
		GenerateURL(gomock.Any()).
		Return("shortURL").
		Times(calls)

	shortener := service.NewShortenerService(repo, mockURLGen, config)

	t.Run("url collisions test", func(t *testing.T) {
		t.Parallel()

		_, err := shortener.ShortenURL("abc")

		require.NoError(t, err)

		start := time.Now()
		_, err = shortener.ShortenURL("cba")

		require.ErrorIs(t, err, e.ErrServiceCollision)
		assert.LessOrEqual(t, time.Since(start), config.URLGenTimeout)
	})
}

func RunStoreRestoreTests(t *testing.T, shortener *service.URLShortener) {
	t.Helper()

	t.Run("ShortenURL", func(t *testing.T) {
		t.Parallel()

		longURL := utils.RandURL()
		shortURL, err := shortener.ShortenURL(longURL)

		require.NoError(t, err)

		restoredURL, err := shortener.GetOriginalURL(shortURL)
		require.NoError(t, err)
		assert.Equal(t, longURL, restoredURL)
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
			GetURL(gomock.Any()).
			Return("", e.ErrTest).
			Times(1),

		mockURLRepo.
			EXPECT().
			GetURL(gomock.Any()).
			Return("", e.ErrRepoNotFound).
			Times(1),

		mockURLRepo.
			EXPECT().
			AddURL(gomock.Any(), gomock.Any()).
			Return(e.ErrTest).
			Times(1),
	)

	shortener := service.NewShortenerService(mockURLRepo, gen, config)

	t.Run("ShortenURLErr", func(t *testing.T) {
		t.Parallel()

		// first call should have problems to get from repo.
		_, err := shortener.ShortenURL("http://localhost:8080")
		require.ErrorIs(t, err, e.ErrServiceInternal)

		// second call should have problems to add to repo.
		_, err = shortener.ShortenURL("http://localhost:8181")
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
			GetURL(gomock.Any()).
			Return("", e.ErrTest).
			Times(1),

		mockURLRepo.
			EXPECT().
			GetURL(gomock.Any()).
			Return("", e.ErrRepoNotFound).
			Times(1),
	)

	shortener := service.NewShortenerService(mockURLRepo, gen, config)
	badURL := utils.RandomString(config.URLsize + 1)

	_, err := shortener.GetOriginalURL(badURL)
	require.ErrorIs(t, err, e.ErrServiceInvalid)

	_, err = shortener.GetOriginalURL("shortURL")
	require.ErrorIs(t, err, e.ErrServiceInternal)

	_, err = shortener.GetOriginalURL("shortURL")
	require.ErrorIs(t, err, e.ErrRepoNotFound)
}

func TestURLShortener(t *testing.T) {
	t.Parallel()

	config := config.DefaultConfig()
	repo := repository.NewInMemoryURLRepository()
	gen := urlgenerator.NewRandURLGenerator(config.URLsize)
	shortener := service.NewShortenerService(repo, gen, config)

	// store and restore good urls.
	for range numStoreRestoreTests {
		RunStoreRestoreTests(t, shortener)
	}

	// colisions and retries.
	RunShortenURLCollisionsTests(t, config, repo)

	// internal errors.
	RunShortenURLErrTests(t, config, gen)
	RunGetOriginalURLErrTests(t, config, gen)
}
