package urlgenerator_test

import (
	"context"
	"math/rand"
	"testing"
	"time"

	"github.com/patraden/ya-practicum-go-shortly/internal/app/domain"
	"github.com/patraden/ya-practicum-go-shortly/internal/app/service/urlgenerator"
)

const (
	charset      = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	numTest      = 20
	shortURLSize = 8
	longURLSize  = 20
)

func RandomString(n int) string {
	seededRand := rand.New(rand.NewSource(time.Now().UnixNano()))
	b := make([]byte, n)

	for i := range b {
		b[i] = charset[seededRand.Intn(len(charset))]
	}

	return string(b)
}

func RunURLGeneratorTests(t *testing.T, generator urlgenerator.URLGenerator, original domain.OriginalURL) {
	t.Helper()
	t.Run("GenerateURL", func(t *testing.T) {
		t.Parallel()

		slug := generator.GenerateSlug(context.Background(), original)

		if !generator.IsValidSlug(slug) {
			t.Errorf("Generated slug invalid: %s", slug.String())
		}
	})
}

func TestRandURLGenerator(t *testing.T) {
	t.Parallel()

	generator := urlgenerator.NewRandURLGenerator(shortURLSize)

	for range numTest {
		original := domain.OriginalURL(RandomString(longURLSize))
		RunURLGeneratorTests(t, generator, original)
	}
}
