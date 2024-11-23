package dto

import "github.com/patraden/ya-practicum-go-shortly/internal/app/domain"

//go:generate easyjson -all
//easyjson:json
type ShortenURLRequest struct {
	LongURL string `json:"url"`
}

//easyjson:json
type ShortenedURLResponse struct {
	ShortURL string `json:"result"`
}

//easyjson:json
type URLPair struct {
	Slug        domain.Slug        `json:"short_url"`
	OriginalURL domain.OriginalURL `json:"original_url"`
}

//easyjson:json
type CorrelatedOriginalURL struct {
	CorrelationID string             `json:"correlation_id"`
	OriginalURL   domain.OriginalURL `json:"original_url"`
}

//easyjson:json
type CorrelatedSlug struct {
	CorrelationID string      `json:"correlation_id"`
	Slug          domain.Slug `json:"short_url"`
}

//easyjson:json
type OriginalURLBatch []CorrelatedOriginalURL

func (b OriginalURLBatch) Originals() []domain.OriginalURL {
	originals := make([]domain.OriginalURL, len(b))
	for i, elem := range b {
		originals[i] = elem.OriginalURL
	}

	return originals
}

//easyjson:json
type SlugBatch []CorrelatedSlug

//easyjson:json
type URLPairBatch []URLPair

func NewURLPairBatch(maps *[]domain.URLMapping, baseURL string) *URLPairBatch {
	res := make(URLPairBatch, len(*maps))

	for i, m := range *maps {
		res[i] = URLPair{
			Slug:        domain.Slug(m.Slug.WithBaseURL(baseURL)),
			OriginalURL: m.OriginalURL,
		}
	}

	return &res
}
