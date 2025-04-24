package domain_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/patraden/ya-practicum-go-shortly/internal/app/domain"
)

func TestSlugString(t *testing.T) {
	t.Parallel()

	s := domain.Slug("abc123")
	assert.Equal(t, "abc123", s.String())
}

func TestSlugWithBaseURL(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		baseURL  string
		slug     domain.Slug
		expected string
	}{
		{"Base URL with trailing slash", "http://short.ly/", "xyz", "http://short.ly/xyz"},
		{"Base URL without trailing slash", "http://short.ly", "xyz", "http://short.ly/xyz"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual := tt.slug.WithBaseURL(tt.baseURL)
			assert.Equal(t, tt.expected, actual)
		})
	}
}

func TestOriginalURLString(t *testing.T) {
	t.Parallel()

	u := domain.OriginalURL("https://example.com")
	assert.Equal(t, "https://example.com", u.String())
}

func TestOriginalURLIsValid(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		url      domain.OriginalURL
		expected bool
	}{
		{"Valid URL", "https://example.com", true},
		{"Valid URL with query", "https://example.com?foo=bar", true},
		{"Invalid URL (missing scheme)", "example.com", false},
		{"Invalid URL (empty string)", "", false},
		{"Invalid URL (only scheme)", "https://", false},
		{"Invalid URL (invalid characters)", "ht@tp://bad_url.com", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, tt.url.IsValid())
		})
	}
}

func TestURLMappingExpiresAfter(t *testing.T) {
	t.Parallel()

	mapping := &domain.URLMapping{CreatedAt: time.Now()}
	duration := time.Hour * 24 * 30 // 30 days

	mapping.ExpiresAfter(duration)
	assert.Equal(t, mapping.CreatedAt.Add(duration), mapping.ExpiresAt)
}

func TestNewURLMapping(t *testing.T) {
	t.Parallel()

	slug := domain.Slug("short123")
	original := domain.OriginalURL("https://example.com")
	userID := domain.NewUserID()

	mapping := domain.NewURLMapping(slug, original, userID)

	require.NotNil(t, mapping)
	assert.Equal(t, slug, mapping.Slug)
	assert.Equal(t, original, mapping.OriginalURL)
	assert.Equal(t, userID, mapping.UserID)
	assert.False(t, mapping.Deleted)
	assert.WithinDuration(t, time.Now().Add(time.Hour*24*730), mapping.ExpiresAt, time.Second*2)
}
