package repository

import (
	"context"
	"sync"

	"github.com/patraden/ya-practicum-go-shortly/internal/app/domain"
	e "github.com/patraden/ya-practicum-go-shortly/internal/app/domain/errors"
	"github.com/patraden/ya-practicum-go-shortly/internal/app/dto"
	"github.com/patraden/ya-practicum-go-shortly/internal/app/memento"
)

// InMemoryURLRepository is an in-memory implementation of the URL repository.
type InMemoryURLRepository struct {
	sync.RWMutex
	values   dto.URLMappings
	uIndex   map[domain.OriginalURL]domain.Slug
	usrIndex map[domain.UserID][]domain.Slug
}

// NewInMemoryURLRepository creates a new InMemoryURLRepository instance.
func NewInMemoryURLRepository() *InMemoryURLRepository {
	return &InMemoryURLRepository{
		RWMutex:  sync.RWMutex{},
		values:   make(dto.URLMappings),
		uIndex:   make(map[domain.OriginalURL]domain.Slug),
		usrIndex: make(map[domain.UserID][]domain.Slug),
	}
}

// AddURLMapping adds a new URL mapping to the repository.
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
	ms.usrIndex[urlMap.UserID] = append(ms.usrIndex[urlMap.UserID], urlMap.Slug)

	return urlMap, nil
}

// GetURLMapping retrieves a URL mapping by its slug.
func (ms *InMemoryURLRepository) GetURLMapping(_ context.Context, slug domain.Slug) (*domain.URLMapping, error) {
	ms.RLock()
	defer ms.RUnlock()

	m, exists := ms.values[slug]
	if !exists {
		return nil, e.ErrSlugNotFound
	}

	return &m, nil
}

// GetUserURLMappings retrieves all URL mappings for a specific user.
func (ms *InMemoryURLRepository) GetUserURLMappings(
	_ context.Context,
	user domain.UserID,
) ([]domain.URLMapping, error) {
	ms.RLock()
	defer ms.RUnlock()

	if _, exists := ms.usrIndex[user]; !exists {
		return nil, e.ErrUserNotFound
	}

	res := make([]domain.URLMapping, len(ms.usrIndex[user]))
	for i, slug := range ms.usrIndex[user] {
		res[i] = ms.values[slug]
	}

	return res, nil
}

// AddURLMappingBatch adds multiple URL mappings in a single batch to the database.
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
		ms.usrIndex[m.UserID] = append(ms.usrIndex[m.UserID], m.Slug)
	}

	return nil
}

// CreateMemento creates a memento of the current state of the repository.
func (ms *InMemoryURLRepository) CreateMemento() (*memento.Memento, error) {
	ms.RLock()
	defer ms.RUnlock()

	cp := dto.URLMappingsCopy(ms.values)

	return memento.NewMemento(cp), nil
}

// RestoreMemento restores the state of the repository from the given memento.
func (ms *InMemoryURLRepository) RestoreMemento(m *memento.Memento) error {
	if m == nil {
		return nil
	}

	ms.Lock()
	defer ms.Unlock()

	// Copy and restore values.
	cp := dto.URLMappingsCopy(m.GetState())
	ms.values = cp

	// Rebuild indexes to maintain consistency with values.
	ms.uIndex = make(map[domain.OriginalURL]domain.Slug)
	ms.usrIndex = make(map[domain.UserID][]domain.Slug)

	for slug, mapping := range ms.values {
		ms.uIndex[mapping.OriginalURL] = slug
		ms.usrIndex[mapping.UserID] = append(ms.usrIndex[mapping.UserID], slug)
	}

	return nil
}

// DelUserURLMappings marks user URL mappings as deleted based on the provided tasks.
func (ms *InMemoryURLRepository) DelUserURLMappings(_ context.Context, tasks []dto.UserSlug) error {
	updateTasks := make([]dto.UserSlug, 0, len(tasks))

	ms.RLock()
	for _, task := range tasks {
		val, ok := ms.values[task.Slug]
		if ok && val.UserID == task.UserID {
			updateTasks = append(updateTasks, task)
		}
	}
	ms.RUnlock()

	ms.Lock()

	for _, task := range updateTasks {
		val := ms.values[task.Slug]
		val.Deleted = true
		ms.values[task.Slug] = val
	}
	ms.Unlock()

	return nil
}

// GetStats retrieves repo statistics.
func (ms *InMemoryURLRepository) GetStats(_ context.Context) (*dto.RepoStats, error) {
	ms.RLock()
	defer ms.RUnlock()

	stats := &dto.RepoStats{
		CountSlugs: int64(len(ms.values)),
		CountUsers: int64(len(ms.usrIndex)),
	}

	return stats, nil
}
