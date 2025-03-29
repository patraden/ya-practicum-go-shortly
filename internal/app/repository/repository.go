package repository

import (
	"context"

	"github.com/patraden/ya-practicum-go-shortly/internal/app/domain"
	"github.com/patraden/ya-practicum-go-shortly/internal/app/dto"
	"github.com/patraden/ya-practicum-go-shortly/internal/app/memento"
)

const errLabel = "repository"

// URLRepository is an interface that defines the methods for interacting with URL mappings in a repository.
type URLRepository interface {
	memento.Originator
	AddURLMapping(ctx context.Context, m *domain.URLMapping) (*domain.URLMapping, error)
	AddURLMappingBatch(ctx context.Context, batch *[]domain.URLMapping) error
	GetURLMapping(ctx context.Context, slug domain.Slug) (*domain.URLMapping, error)
	GetUserURLMappings(ctx context.Context, user domain.UserID) ([]domain.URLMapping, error)
	DelUserURLMappings(ctx context.Context, tasks []dto.UserSlug) error
	GetStats(ctx context.Context) (*dto.RepoStats, error)
}
