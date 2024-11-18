package repository

import (
	"context"
	"sync"

	"github.com/patraden/ya-practicum-go-shortly/internal/app/domain"
	e "github.com/patraden/ya-practicum-go-shortly/internal/app/domain/errors"
	"github.com/patraden/ya-practicum-go-shortly/internal/app/dto"
	"github.com/patraden/ya-practicum-go-shortly/internal/app/memento"
)

type InMemoryURLRepository struct {
	sync.RWMutex
	values dto.URLMappings
	uIndex map[domain.OriginalURL]domain.Slug
}

func NewInMemoryURLRepository() *InMemoryURLRepository {
	return &InMemoryURLRepository{
		RWMutex: sync.RWMutex{},
		values:  make(dto.URLMappings),
		uIndex:  make(map[domain.OriginalURL]domain.Slug),
	}
}

func (ms *InMemoryURLRepository) AddURLMapping(
	_ context.Context,
	urlMap *domain.URLMapping,
) (*domain.URLMapping, error) {
	ms.Lock()
	defer ms.Unlock()

	if _, exists := ms.values[urlMap.Slug]; exists {
		return urlMap, e.ErrSlugExists
	}

	if _, exists := ms.uIndex[urlMap.OriginalURL]; exists {
		slug := ms.uIndex[urlMap.OriginalURL]
		urlMapping := ms.values[slug]

		return &urlMapping, e.ErrOriginalExists
	}

	ms.values[urlMap.Slug] = *urlMap
	ms.uIndex[urlMap.OriginalURL] = urlMap.Slug

	return urlMap, nil
}

func (ms *InMemoryURLRepository) GetURLMapping(_ context.Context, slug domain.Slug) (*domain.URLMapping, error) {
	ms.RLock()
	defer ms.RUnlock()

	m, exists := ms.values[slug]
	if !exists {
		return nil, e.ErrSlugNotFound
	}

	return &m, nil
}

func (ms *InMemoryURLRepository) AddURLMappingBatch(_ context.Context, batch *[]domain.URLMapping) error {
	ms.Lock()
	defer ms.Unlock()

	// Validate entire batch to simulate transactional behavior.
	for _, m := range *batch {
		if _, exists := ms.values[m.Slug]; exists {
			return e.ErrSlugExists
		}

		if _, exists := ms.uIndex[m.OriginalURL]; exists {
			return e.ErrOriginalExists
		}
	}

	// No conflicts found; proceed with adding to maps.
	for _, m := range *batch {
		ms.values[m.Slug] = m
		ms.uIndex[m.OriginalURL] = m.Slug
	}

	return nil
}

func (ms *InMemoryURLRepository) CreateMemento() (*memento.Memento, error) {
	ms.RLock()
	defer ms.RUnlock()

	cp := dto.URLMappingsCopy(ms.values)

	return memento.NewMemento(cp), nil
}

func (ms *InMemoryURLRepository) RestoreMemento(m *memento.Memento) error {
	if m == nil {
		return nil
	}

	ms.Lock()
	defer ms.Unlock()

	// Copy and restore values.
	cp := dto.URLMappingsCopy(m.GetState())
	ms.values = cp

	// Rebuild index to maintain consistency with values.
	ms.uIndex = make(map[domain.OriginalURL]domain.Slug)
	for slug, mapping := range ms.values {
		ms.uIndex[mapping.OriginalURL] = slug
	}

	return nil
}
