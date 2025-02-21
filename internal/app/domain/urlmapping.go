package domain

import (
	"net/url"
	"time"
)

const (
	defaultExpiration = time.Hour * 24 * 730 // 2 years
	errLabel          = "domain"
)

// Slug represents a short, unique identifier for a URL.
type Slug string

// OriginalURL represents the original, long-form URL before shortening.
type OriginalURL string

// String returns the string representation of the Slug.
func (s Slug) String() string {
	return string(s)
}

// WithBaseURL constructs a full URL by appending the Slug to a given base URL.
func (s *Slug) WithBaseURL(baseURL string) string {
	if baseURL[len(baseURL)-1] != '/' {
		baseURL += "/"
	}

	return baseURL + s.String()
}

// String returns the string representation of the OriginalURL.
func (u OriginalURL) String() string {
	return string(u)
}

// IsValid checks whether the OriginalURL is a valid URL.
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

// URLMapping represents a mapping between a shortened URL (Slug) and its OriginalURL.
type URLMapping struct {
	Slug        Slug        `json:"short_url"`
	OriginalURL OriginalURL `json:"original_url"`
	UserID      UserID      `json:"user_id"`
	CreatedAt   time.Time   `json:"created_at"`
	ExpiresAt   time.Time   `json:"expires_at"`
	Deleted     bool        `json:"is_deleted"`
}

// ExpiresAfter sets the expiration time of the URLMapping based on a given duration.
func (m *URLMapping) ExpiresAfter(duration time.Duration) {
	m.ExpiresAt = m.CreatedAt.Add(duration)
}

// NewURLMapping creates a new URLMapping instance with the given Slug, OriginalURL, and UserID.
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
