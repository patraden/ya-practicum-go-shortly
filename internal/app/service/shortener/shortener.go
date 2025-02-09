package shortener

import (
	"context"

	"github.com/patraden/ya-practicum-go-shortly/internal/app/domain"
	"github.com/patraden/ya-practicum-go-shortly/internal/app/dto"
)

const errLabel = "shortener"

type URLShortener interface {
	ShortenURL(ctx context.Context, original domain.OriginalURL) (domain.Slug, error)
	ShortenURLBatch(ctx context.Context, batch *dto.OriginalURLBatch) (*dto.SlugBatch, error)
	GetOriginalURL(ctx context.Context, slug domain.Slug) (domain.OriginalURL, error)
	GetUserURLs(ctx context.Context) (*dto.URLPairBatch, error)
}
