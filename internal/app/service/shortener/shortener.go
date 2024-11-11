package shortener

import (
	"context"

	"github.com/patraden/ya-practicum-go-shortly/internal/app/domain"
	"github.com/patraden/ya-practicum-go-shortly/internal/app/dto"
)

type URLShortener interface {
	ShortenURL(ctx context.Context, original domain.OriginalURL) (*domain.URLMapping, error)
	GetOriginalURL(ctx context.Context, slug domain.Slug) (*domain.URLMapping, error)
	ShortenURLBatch(ctx context.Context, batch *dto.OriginalURLBatch) (dto.SlugBatch, error)
}
