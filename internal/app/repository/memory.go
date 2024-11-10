package repository

import (
	"context"
	"sync"

	"github.com/patraden/ya-practicum-go-shortly/internal/app/domain"
	"github.com/patraden/ya-practicum-go-shortly/internal/app/dto"
	e "github.com/patraden/ya-practicum-go-shortly/internal/app/errors"
	"github.com/patraden/ya-practicum-go-shortly/internal/app/memento"
)

type InMemoryURLRepository struct {
	sync.RWMutex
	values dto.URLMappings
}

func NewInMemoryURLRepository() *InMemoryURLRepository {
	return &InMemoryURLRepository{
		RWMutex: sync.RWMutex{},
		values:  make(dto.URLMappings),
	}
}

func (ms *InMemoryURLRepository) AddURLMapping(_ context.Context, m *domain.URLMapping) error {
	ms.Lock()
	defer ms.Unlock()

	if _, ok := ms.values[m.Slug]; ok {
		return e.ErrRepoExists
	}

	ms.values[m.Slug] = *m

	return nil
}

func (ms *InMemoryURLRepository) GetURLMapping(_ context.Context, slug domain.Slug) (*domain.URLMapping, error) {
	ms.RLock()
	defer ms.RUnlock()

	m, ok := ms.values[slug]
	if !ok {
		return nil, e.ErrRepoNotFound
	}

	return &m, nil
}

func (ms *InMemoryURLRepository) CreateMemento() (*memento.Memento, error) {
	ms.RLock()
	defer ms.RUnlock()

	cp := dto.URLMappingsCopy(ms.values)

	return memento.NewMemento(cp), nil
}

func (ms *InMemoryURLRepository) RestoreMemento(m *memento.Memento) error {
	ms.Lock()
	defer ms.Unlock()

	cp := dto.URLMappingsCopy(m.GetState())
	ms.values = cp

	return nil
}
