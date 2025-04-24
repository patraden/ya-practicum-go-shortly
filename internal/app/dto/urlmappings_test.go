package dto_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/patraden/ya-practicum-go-shortly/internal/app/domain"
	"github.com/patraden/ya-practicum-go-shortly/internal/app/dto"
)

func TestURLMappingsCopy(t *testing.T) {
	t.Parallel()

	slug1 := domain.Slug("slug1")
	urlMap1 := *domain.NewURLMapping(
		slug1,
		domain.OriginalURL("url1"),
		domain.NewUserID(),
	)

	slug2 := domain.Slug("slug2")
	urlMap2 := *domain.NewURLMapping(
		slug2,
		domain.OriginalURL("url2"),
		domain.NewUserID(),
	)

	t.Run("Copy test", func(t *testing.T) {
		originalMap := dto.URLMappings{
			slug1: urlMap1,
			slug2: urlMap2,
		}

		copiedMap := dto.URLMappingsCopy(originalMap)

		assert.Equal(t, originalMap, copiedMap)

		// Ensure the copied map is a deep copy (modifying one should not affect the other)
		copiedMap[slug1] = urlMap2
		assert.NotEqual(t, originalMap[slug1], copiedMap[slug1])
	})
}
