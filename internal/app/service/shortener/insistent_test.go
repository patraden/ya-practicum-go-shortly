package shortener_test

import (
	"context"
	"testing"
	"time"

	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	"github.com/patraden/ya-practicum-go-shortly/internal/app/config"
	"github.com/patraden/ya-practicum-go-shortly/internal/app/domain"
	e "github.com/patraden/ya-practicum-go-shortly/internal/app/domain/errors"
	"github.com/patraden/ya-practicum-go-shortly/internal/app/dto"
	"github.com/patraden/ya-practicum-go-shortly/internal/app/mock"
	"github.com/patraden/ya-practicum-go-shortly/internal/app/service/shortener"
)

func TestShortenURL(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := mock.NewMockURLRepository(ctrl)
	urlGen := mock.NewMockURLGenerator(ctrl)
	config := config.DefaultConfig()
	log := zerolog.New(nil)
	svc := shortener.NewInsistentShortener(repo, urlGen, config, &log)
	ctx := context.Background()

	t.Run("successfully shortens a URL", func(t *testing.T) {
		urlMapping := domain.NewURLMapping("slug1", "http://example.com")

		urlGen.EXPECT().GenerateSlug(gomock.Any(), urlMapping.OriginalURL).Return(urlMapping.Slug).Times(1)
		repo.EXPECT().GetURLMapping(gomock.Any(), urlMapping.Slug).Return(nil, e.ErrSlugNotFound)
		repo.EXPECT().AddURLMapping(gomock.Any(), gomock.Any()).Return(urlMapping, nil)

		result, err := svc.ShortenURL(ctx, urlMapping.OriginalURL)
		require.NoError(t, err)
		assert.Equal(t, urlMapping.Slug, result)
	})

	t.Run("returns error on slug collision", func(t *testing.T) {
		original := domain.OriginalURL("http://example.com")
		slug := domain.Slug("slug1")
		urlMapping := domain.NewURLMapping("slug1", "http://example.com")
		// we expected calles as per linear backoff
		calls := int(config.URLGenTimeout / config.URLGenRetryInterval)

		urlGen.EXPECT().GenerateSlug(gomock.Any(), original).Return(slug).Times(calls)
		repo.EXPECT().GetURLMapping(gomock.Any(), slug).Return(urlMapping, nil).Times(calls)

		start := time.Now()
		result, err := svc.ShortenURL(ctx, original)
		// ensure retries are not exceeding max time
		assert.LessOrEqual(t, time.Since(start), config.URLGenTimeout)
		// ensure colission error is returned
		require.ErrorIs(t, err, e.ErrSlugCollision)
		assert.Equal(t, domain.Slug(""), result)
	})

	t.Run("returns internal error on unexpected repository failure", func(t *testing.T) {
		original := domain.OriginalURL("http://example.com")
		slug := domain.Slug("short1")

		urlGen.EXPECT().GenerateSlug(gomock.Any(), original).Return(slug)
		repo.EXPECT().GetURLMapping(gomock.Any(), slug).Return(nil, e.ErrTestGeneral)

		result, err := svc.ShortenURL(context.Background(), original)
		require.ErrorIs(t, err, e.ErrShortenerInternal)
		assert.Equal(t, domain.Slug(""), result)
	})

	t.Run("returns error if original URL already exists", func(t *testing.T) {
		original := domain.OriginalURL("http://example.com")
		slugDup := domain.Slug("slug2")
		urlMapping := domain.NewURLMapping("slug1", original)

		urlGen.EXPECT().GenerateSlug(gomock.Any(), original).Return(slugDup).Times(1)
		repo.EXPECT().GetURLMapping(gomock.Any(), slugDup).Return(nil, e.ErrSlugNotFound).Times(1)
		repo.EXPECT().AddURLMapping(gomock.Any(), gomock.Any()).Return(urlMapping, e.ErrOriginalExists).Times(1)

		result, err := svc.ShortenURL(ctx, original)
		require.ErrorIs(t, err, e.ErrOriginalExists)
		assert.Equal(t, "slug1", result.String())
	})
}

func TestGetOriginalURL(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := mock.NewMockURLRepository(ctrl)
	urlGen := mock.NewMockURLGenerator(ctrl)
	config := config.DefaultConfig()
	log := zerolog.New(nil)
	svc := shortener.NewInsistentShortener(repo, urlGen, config, &log)
	ctx := context.Background()

	t.Run("successfully retrieves original URL", func(t *testing.T) {
		slug := domain.Slug("short1")
		original := domain.OriginalURL("http://example.com")
		urlMapping := domain.NewURLMapping("short1", "http://example.com")

		urlGen.EXPECT().IsValidSlug(slug).Return(true)
		repo.EXPECT().GetURLMapping(gomock.Any(), slug).Return(urlMapping, nil)

		result, err := svc.GetOriginalURL(ctx, slug)
		require.NoError(t, err)
		assert.Equal(t, original.String(), result.String())
	})

	t.Run("returns not found error for unknown slug", func(t *testing.T) {
		slug := domain.Slug("unknown")

		urlGen.EXPECT().IsValidSlug(slug).Return(true)
		repo.EXPECT().GetURLMapping(gomock.Any(), slug).Return(nil, e.ErrSlugNotFound)

		result, err := svc.GetOriginalURL(ctx, slug)
		require.ErrorIs(t, err, e.ErrSlugNotFound)
		assert.Equal(t, domain.OriginalURL(""), result)
	})

	t.Run("returns internal error on unexpected failure", func(t *testing.T) {
		slug := domain.Slug("short1")

		urlGen.EXPECT().IsValidSlug(slug).Return(true)
		repo.EXPECT().GetURLMapping(gomock.Any(), slug).Return(nil, e.ErrTestGeneral)

		result, err := svc.GetOriginalURL(ctx, slug)
		require.ErrorIs(t, err, e.ErrShortenerInternal)
		assert.Equal(t, domain.OriginalURL(""), result)
	})
}

func TestShortenURLBatch(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := mock.NewMockURLRepository(ctrl)
	urlGen := mock.NewMockURLGenerator(ctrl)
	config := config.DefaultConfig()
	log := zerolog.New(nil)
	svc := shortener.NewInsistentShortener(repo, urlGen, config, &log)
	ctx := context.Background()

	t.Run("successfully shortens a batch of URLs", func(t *testing.T) {
		originals := dto.OriginalURLBatch{
			{CorrelationID: "1", OriginalURL: "http://example1.com"},
			{CorrelationID: "2", OriginalURL: "http://example2.com"},
		}

		slugs := []domain.Slug{"short1", "short2"}
		expected := dto.SlugBatch{
			{CorrelationID: "1", Slug: slugs[0]},
			{CorrelationID: "2", Slug: slugs[1]},
		}

		urlGen.EXPECT().GenerateSlugs(gomock.Any(), originals.Originals()).Return(slugs, nil).Times(1)
		repo.EXPECT().AddURLMappingBatch(gomock.Any(), gomock.Any()).Return(nil).Times(1)

		result, err := svc.ShortenURLBatch(ctx, &originals)
		require.NoError(t, err)
		assert.Equal(t, expected, *result)
	})

	t.Run("returns slug collision error for batch when slugs exist", func(t *testing.T) {
		orig := dto.OriginalURLBatch{{CorrelationID: "1", OriginalURL: "http://example1.com"}}
		slugs := []domain.Slug{"short1"}
		// we expected calles as per linear backoff
		calls := int(config.URLGenTimeout / config.URLGenRetryInterval)

		urlGen.EXPECT().GenerateSlugs(gomock.Any(), orig.Originals()).Return(slugs, nil).Times(calls)
		repo.EXPECT().AddURLMappingBatch(gomock.Any(), gomock.Any()).Return(e.ErrSlugExists).Times(calls)

		start := time.Now()
		result, err := svc.ShortenURLBatch(ctx, &orig)
		// ensure retries are not exceeding max time
		assert.LessOrEqual(t, time.Since(start), config.URLGenTimeout)
		require.ErrorIs(t, err, e.ErrSlugCollision)
		assert.Nil(t, result)
	})

	t.Run("returns internal error on batch processing failure", func(t *testing.T) {
		orig := dto.OriginalURLBatch{{CorrelationID: "1", OriginalURL: "http://example1.com"}}
		slugs := []domain.Slug{"short1"}

		urlGen.EXPECT().GenerateSlugs(gomock.Any(), orig.Originals()).Return(slugs, nil)
		repo.EXPECT().AddURLMappingBatch(gomock.Any(), gomock.Any()).Return(e.ErrTestGeneral)

		_, err := svc.ShortenURLBatch(context.Background(), &orig)
		require.ErrorIs(t, err, e.ErrShortenerInternal)
	})
}
