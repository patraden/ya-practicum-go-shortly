package urlgenerator

import (
	"context"

	"github.com/patraden/ya-practicum-go-shortly/internal/app/domain"
)

// Short URL generator app service interface
// Decided to dedicate an interface for this service
// as potentially throughout the course of development
// there might be different implemnetations
// like random, incremental, hash based etc.
type URLGenerator interface {
	GenerateSlug(ctx context.Context, original domain.OriginalURL) domain.Slug
	GenerateSlugs(ctx context.Context, originals []domain.OriginalURL) ([]domain.Slug, error)
	IsValidSlug(slug domain.Slug) bool
}
