package urlgenerator

import (
	"context"
	"fmt"
	"math/rand"
	"regexp"

	"github.com/patraden/ya-practicum-go-shortly/internal/app/domain"
	e "github.com/patraden/ya-practicum-go-shortly/internal/app/domain/errors"
)

// RandURLGenerator generates random slugs of a fixed length.
type RandURLGenerator struct {
	length int
}

// NewRandURLGenerator creates a new instance of RandURLGenerator with a specified length.
func NewRandURLGenerator(l int) *RandURLGenerator {
	return &RandURLGenerator{
		length: l,
	}
}

// GenerateSlug generates a random slug of a fixed length for a given original URL.
// It selects characters from a predefined set (uppercase letters, lowercase letters, and digits)
// and returns the generated slug.
func (g *RandURLGenerator) GenerateSlug(_ context.Context, _ domain.OriginalURL) domain.Slug {
	charSets := []string{
		"ABCDEFGHIJKLMNOPQRSTUVWXYZ", // A-Z
		"abcdefghijklmnopqrstuvwxyz", // a-z
		"0123456789",                 // 0-9
	}

	bytes := make([]byte, g.length)
	for i := range bytes {
		charSet := charSets[rand.Intn(len(charSets))]
		bytes[i] = charSet[rand.Intn(len(charSet))]
	}

	return domain.Slug(string(bytes))
}

// IsValidSlug checks if the given slug matches the required format: a string of letters and digits
// of the specified length.
func (g *RandURLGenerator) IsValidSlug(slug domain.Slug) bool {
	regexPattern := fmt.Sprintf(`^/?[a-zA-Z0-9]{%d}$`, g.length)
	validShortURL := regexp.MustCompile(regexPattern)

	return validShortURL.MatchString(slug.String())
}

// GenerateSlugs generates unique random slugs for a batch of original URLs.
// It ensures that all generated slugs are unique by using a map to track already generated slugs.
func (g *RandURLGenerator) GenerateSlugs(ctx context.Context, originals []domain.OriginalURL) ([]domain.Slug, error) {
	unique := make(map[domain.Slug]interface{}, len(originals))
	res := make([]domain.Slug, len(originals))
	i := 0

	for i < len(originals) {
		select {
		case <-ctx.Done():
			return []domain.Slug{}, e.ErrURLGenGenerateSlug
		default:
			slug := g.GenerateSlug(ctx, originals[i])
			if _, exists := unique[slug]; !exists {
				res[i] = slug
				i++
			}
		}
	}

	return res, nil
}
