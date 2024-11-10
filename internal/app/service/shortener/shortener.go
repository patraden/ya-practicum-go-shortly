package shortener

import (
	"context"

	"github.com/patraden/ya-practicum-go-shortly/internal/app/domain"
)

type URLShortener interface {
	ShortenURL(ctx context.Context, original domain.OriginalURL) (*domain.URLMapping, error)
	GetOriginalURL(ctx context.Context, slug domain.Slug) (*domain.URLMapping, error)
}
