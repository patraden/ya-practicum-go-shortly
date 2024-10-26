package urlgenerator_test

import (
	"math/rand"
	"testing"
	"time"

	"github.com/patraden/ya-practicum-go-shortly/internal/app/urlgenerator"
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

func RunURLGeneratorTests(t *testing.T, generator urlgenerator.URLGenerator, url string) {
	t.Helper()
	t.Run("GenerateURL", func(t *testing.T) {
		t.Parallel()

		shortURL := generator.GenerateURL(url)

		if !generator.IsValidURL(shortURL) {
			t.Errorf("Generated URL is not valid: %s", shortURL)
		}
	})
}

func TestRandURLGenerator(t *testing.T) {
	t.Parallel()

	generator := urlgenerator.NewRandURLGenerator(shortURLSize)

	for range numTest {
		RunURLGeneratorTests(t, generator, RandomString(longURLSize))
	}
}
