package domain

import (
	"net/url"
	"time"
)

const defaultExpiration = time.Hour * 24 * 730 // 2 years

type (
	Slug        string
	OriginalURL string
)

func (s Slug) String() string {
	return string(s)
}

func (s *Slug) WithBaseURL(baseURL string) string {
	if baseURL[len(baseURL)-1] != '/' {
		baseURL += "/"
	}

	return baseURL + s.String()
}

func (u OriginalURL) String() string {
	return string(u)
}

func (u OriginalURL) IsValid() bool {
	parsedURL, err := url.ParseRequestURI(u.String())
	if err != nil {
		return false
	}

	if parsedURL.Scheme == "" || parsedURL.Host == "" {
		return false
	}

	return true
}

type URLMapping struct {
	Slug        Slug        `json:"short_url"`
	OriginalURL OriginalURL `json:"original_url"`
	CreatedAt   time.Time   `json:"created_at"`
	ExpiresAt   time.Time   `json:"expires_at"`
}

func (m *URLMapping) ExpiresAfter(duration time.Duration) {
	m.ExpiresAt = m.CreatedAt.Add(duration)
}

func NewURLMapping(slug Slug, original OriginalURL) *URLMapping {
	m := &URLMapping{
		Slug:        slug,
		OriginalURL: original,
		CreatedAt:   time.Now(),
	}

	m.ExpiresAfter(defaultExpiration)

	return m
}
