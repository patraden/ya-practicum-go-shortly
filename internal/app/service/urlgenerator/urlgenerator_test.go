package urlgenerator_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/patraden/ya-practicum-go-shortly/internal/app/domain"
	"github.com/patraden/ya-practicum-go-shortly/internal/app/service/urlgenerator"
	"github.com/patraden/ya-practicum-go-shortly/internal/app/utils"
)

const (
	numTest      = 20
	shortURLSize = 8
	longURLSize  = 20
)

func TestRandURLGenerator(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	generator := urlgenerator.NewRandURLGenerator(shortURLSize)
	originals := []domain.OriginalURL{}

	for range numTest {
		original := domain.OriginalURL(utils.RandString(longURLSize))
		originals = append(originals, original)

		t.Run("GenerateURL", func(t *testing.T) {
			t.Parallel()

			slug := generator.GenerateSlug(ctx, original)

			if !generator.IsValidSlug(slug) {
				t.Errorf("Generated slug invalid: %s", slug.String())
			}
		})
	}

	t.Run("GenerateURLs", func(t *testing.T) {
		slugs, err := generator.GenerateSlugs(ctx, originals)
		require.NoError(t, err)

		assert.Len(t, slugs, numTest)

		for _, slug := range slugs {
			if !generator.IsValidSlug(slug) {
				t.Errorf("Generated slug invalid: %s", slug.String())
			}
		}
	})
}
