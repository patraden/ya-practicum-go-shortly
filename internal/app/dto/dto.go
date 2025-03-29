package dto

import "github.com/patraden/ya-practicum-go-shortly/internal/app/domain"

// ShortenURLRequest represents a request payload to shorten a URL.
//
//easyjson:json
type ShortenURLRequest struct {
	LongURL string `json:"url"` // The original long URL to be shortened.
}

// ShortenedURLResponse represents the response containing a shortened URL.
//
//easyjson:json
type ShortenedURLResponse struct {
	ShortURL string `json:"result"` // The generated short URL.
}

// URLPair represents a mapping between a short slug and its corresponding original URL.
//
//easyjson:json
type URLPair struct {
	Slug        domain.Slug        `json:"short_url"`    // The shortened slug.
	OriginalURL domain.OriginalURL `json:"original_url"` // The original full URL.
}

// CorrelatedOriginalURL represents an original URL with a correlation ID.
//
// It is used for batch processing when a correlation ID needs to be associated with each URL.
//
//easyjson:json
type CorrelatedOriginalURL struct {
	CorrelationID string             `json:"correlation_id"` // A unique identifier for correlation.
	OriginalURL   domain.OriginalURL `json:"original_url"`   // The original URL.
}

// CorrelatedSlug represents a shortened URL (slug) with a correlation ID.
//
//easyjson:json
type CorrelatedSlug struct {
	CorrelationID string      `json:"correlation_id"` // A unique identifier for correlation.
	Slug          domain.Slug `json:"short_url"`      // The generated short URL (slug).
}

// UserSlug represents a mapping between a user and their shortened URL slug.
type UserSlug struct {
	Slug   domain.Slug   // The shortened slug.
	UserID domain.UserID // The user who owns the slug.
}

// OriginalURLBatch is a batch of correlated original URLs.
//
//easyjson:json
type OriginalURLBatch []CorrelatedOriginalURL

// Originals extracts and returns a slice of original URLs from the batch.
func (b OriginalURLBatch) Originals() []domain.OriginalURL {
	originals := make([]domain.OriginalURL, len(b))
	for i, elem := range b {
		originals[i] = elem.OriginalURL
	}

	return originals
}

// SlugBatch represents a batch of correlated shortened URLs.
//
//easyjson:json
type SlugBatch []CorrelatedSlug

// UserSlugBatch represents a batch of slugs belonging to users.
//
//easyjson:json
type UserSlugBatch []domain.Slug

// URLPairBatch represents a batch of URL pairs (shortened and original).
//
//easyjson:json
type URLPairBatch []URLPair

// NewURLPairBatch creates a new batch of URL pairs from a slice of URL mappings.
//
// It converts each mapping into a URLPair, appending the base URL to the shortened slug.
func NewURLPairBatch(maps *[]domain.URLMapping, baseURL string) *URLPairBatch {
	res := make(URLPairBatch, len(*maps))

	for i, m := range *maps {
		res[i] = URLPair{
			Slug:        domain.Slug(m.Slug.WithBaseURL(baseURL)), // Append base URL to slug.
			OriginalURL: m.OriginalURL,
		}
	}

	return &res
}

// URLStats represents stats request response content.
//
//easyjson:json
type RepoStats struct {
	CountSlugs int64 `json:"urls"`
	CountUsers int64 `json:"users"`
}
