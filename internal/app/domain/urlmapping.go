package domain

import (
	"net/url"
	"time"
)

const (
	defaultExpiration = time.Hour * 24 * 730 // 2 years
	errLabel          = "domain"
)

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
	UserID      UserID      `json:"user_id"`
	CreatedAt   time.Time   `json:"created_at"`
	ExpiresAt   time.Time   `json:"expires_at"`
	Deleted     bool        `json:"is_deleted"`
}

func (m *URLMapping) ExpiresAfter(duration time.Duration) {
	m.ExpiresAt = m.CreatedAt.Add(duration)
}

func NewURLMapping(slug Slug, original OriginalURL, userID UserID) *URLMapping {
	m := &URLMapping{
		Slug:        slug,
		OriginalURL: original,
		UserID:      userID,
		CreatedAt:   time.Now(),
		Deleted:     false,
	}

	m.ExpiresAfter(defaultExpiration)

	return m
}
