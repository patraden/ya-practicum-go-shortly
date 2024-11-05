package repository

import (
	"sync"

	e "github.com/patraden/ya-practicum-go-shortly/internal/app/errors"
)

type InMemoryURLRepository struct {
	sync.RWMutex
	values map[string]string
}

func NewInMemoryURLRepository() *InMemoryURLRepository {
	return &InMemoryURLRepository{
		RWMutex: sync.RWMutex{},
		values:  map[string]string{},
	}
}

func (ms *InMemoryURLRepository) AddURL(shortURL string, longURL string) error {
	ms.Lock()
	defer ms.Unlock()

	if _, ok := ms.values[shortURL]; ok {
		return e.ErrRepoExists
	}

	ms.values[shortURL] = longURL

	return nil
}

func (ms *InMemoryURLRepository) GetURL(shortURL string) (string, error) {
	ms.RLock()
	defer ms.RUnlock()

	value, ok := ms.values[shortURL]
	if !ok {
		return "", e.ErrRepoNotFound
	}

	return value, nil
}

func (ms *InMemoryURLRepository) DelURL(shortURL string) error {
	ms.Lock()
	defer ms.Unlock()
	delete(ms.values, shortURL)

	return nil
}

func (ms *InMemoryURLRepository) deepCopyValues() map[string]string {
	cp := make(map[string]string)
	for k, v := range ms.values {
		cp[k] = v
	}

	return cp
}

func (ms *InMemoryURLRepository) CreateMemento() (*Memento, error) {
	ms.RLock()
	defer ms.RUnlock()

	return NewURLRepositoryState(ms.deepCopyValues()), nil
}

func (ms *InMemoryURLRepository) RestoreMemento(m *Memento) error {
	ms.Lock()
	defer ms.Unlock()

	cp := make(map[string]string)
	for k, v := range m.GetState() {
		cp[k] = v
	}

	ms.values = cp

	return nil
}
