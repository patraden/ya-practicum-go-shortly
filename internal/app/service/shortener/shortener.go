package shortener

import (
	"context"

	"github.com/patraden/ya-practicum-go-shortly/internal/app/domain"
	"github.com/patraden/ya-practicum-go-shortly/internal/app/dto"
)

const errLabel = "shortener"

// URLShortener defines the interface for a URL shortener service.
// It includes methods for shortening individual URLs, handling batches of URLs,
// retrieving original URLs by slug, and fetching a user's URL mappings.
type URLShortener interface {
	ShortenURL(ctx context.Context, original domain.OriginalURL) (domain.Slug, error)
	ShortenURLBatch(ctx context.Context, batch *dto.OriginalURLBatch) (*dto.SlugBatch, error)
	GetOriginalURL(ctx context.Context, slug domain.Slug) (domain.OriginalURL, error)
	GetUserURLs(ctx context.Context) (*dto.URLPairBatch, error)
}
