package urlgenerator

import (
	"context"
	"fmt"
	"math/rand"
	"regexp"

	"github.com/patraden/ya-practicum-go-shortly/internal/app/domain"
)

type RandURLGenerator struct {
	length int
}

func NewRandURLGenerator(l int) *RandURLGenerator {
	return &RandURLGenerator{
		length: l,
	}
}

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

func (g *RandURLGenerator) IsValidSlug(slug domain.Slug) bool {
	regexPattern := fmt.Sprintf(`^/?[a-zA-Z0-9]{%d}$`, g.length)
	validShortURL := regexp.MustCompile(regexPattern)

	return validShortURL.MatchString(slug.String())
}
