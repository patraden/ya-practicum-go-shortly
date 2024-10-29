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
	ms.Lock()
	defer ms.Unlock()

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

func (ms *InMemoryURLRepository) CreateMemento() (*Memento, error) {
	ms.Lock()
	defer ms.Unlock()

	return NewURLRepositoryState(ms.values), nil
}

func (ms *InMemoryURLRepository) RestoreMemento(m *Memento) error {
	ms.Lock()
	defer ms.Unlock()
	ms.values = m.GetState()

	return nil
}
