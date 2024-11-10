package repository

import (
	"context"

	"github.com/patraden/ya-practicum-go-shortly/internal/app/domain"
)

type URLRepository interface {
	AddURLMapping(ctx context.Context, m *domain.URLMapping) error
	GetURLMapping(ctx context.Context, slug domain.Slug) (*domain.URLMapping, error)
}
