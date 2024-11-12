package repository

import (
	"context"

	"github.com/patraden/ya-practicum-go-shortly/internal/app/domain"
)

const errLabel = "repository"

type URLRepository interface {
	AddURLMapping(ctx context.Context, m *domain.URLMapping) (*domain.URLMapping, error)
	GetURLMapping(ctx context.Context, slug domain.Slug) (*domain.URLMapping, error)
	AddURLMappingBatch(ctx context.Context, batch *[]domain.URLMapping) error
}
