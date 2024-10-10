package repository

import (
	"sync"

	e "github.com/patraden/ya-practicum-go-shortly/internal/app/errors"
)

type MapURLRepository struct {
	sync.RWMutex
	URLRepository
	values map[string]string
}

func NewMapURLRepository() *MapURLRepository {
	return &MapURLRepository{
		values: map[string]string{},
	}
}

func (ms *MapURLRepository) AddURL(shortURL string, longURL string) (string, error) {
	ms.Lock()
	defer ms.Unlock()
	_, ok := ms.values[shortURL]
	if ok {
		return "", e.ErrExists
	}
	ms.values[shortURL] = longURL
	return shortURL, nil
}

func (ms *MapURLRepository) GetURL(shortURL string) (string, error) {
	ms.Lock()
	defer ms.Unlock()

	value, ok := ms.values[shortURL]
	if !ok {
		return "", e.ErrNotFound
	}
	return value, nil
}

func (ms *MapURLRepository) DelURL(shortURL string) error {
	ms.Lock()
	defer ms.Unlock()
	delete(ms.values, shortURL)
	return nil
}
