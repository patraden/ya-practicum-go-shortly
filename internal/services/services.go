package services

import (
	"fmt"

	"github.com/patraden/ya-practicum-go-shortly/internal/storage"
	"github.com/patraden/ya-practicum-go-shortly/internal/urlgen"
)

type LinkStore struct {
	maxGenAttemps int
	urlGenerator  urlgen.URLGenerator
	storage       storage.BasicKVStorage
}

func NewSimpleLinkStore() *LinkStore {
	return &LinkStore{
		maxGenAttemps: 1000,
		urlGenerator:  urlgen.NewRandURLGenerator(8),
		storage:       storage.NewMapStorage(),
	}
}

func (ls *LinkStore) Store(longURL string) (string, error) {
	attemps := 0
	shortURL := ls.urlGenerator.GenerateShortURL(longURL)
	_, err := ls.storage.Get(shortURL)

	for err == nil && attemps < ls.maxGenAttemps {
		shortURL := ls.urlGenerator.GenerateShortURL(longURL)
		_, err = ls.storage.Get(shortURL)
		attemps += 1
	}

	if attemps == ls.maxGenAttemps && err == nil {
		return "", fmt.Errorf("out of shortlinks")
	}

	return ls.storage.Add(shortURL, longURL)
}

func (ls *LinkStore) ReStore(shortURL string) (string, error) {
	return ls.storage.Get(shortURL)
}
