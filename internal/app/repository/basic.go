package repository

import (
	"fmt"

	"github.com/patraden/ya-practicum-go-shortly/internal/app/storage"
	"github.com/patraden/ya-practicum-go-shortly/internal/app/urlgenerator"
)

type BasicLinkRepository struct {
	LinkRepository
	maxGenAttemps int
	urlGenerator  urlgenerator.URLGenerator
	storage       storage.BasicKVStorage
}

func NewBasicLinkRepository() *BasicLinkRepository {
	return &BasicLinkRepository{
		maxGenAttemps: 1000,
		urlGenerator:  urlgenerator.NewRandURLGenerator(8),
		storage:       storage.NewMapStorage(),
	}
}

func (lr *BasicLinkRepository) Store(longURL string) (string, error) {
	shortURL, err := lr.urlGenerator.GenerateURL(longURL)
	if err != nil {
		return "", err
	}

	attemps := 0
	_, err = lr.storage.Get(shortURL)
	for err == nil && attemps < lr.maxGenAttemps {
		shortURL, _ := lr.urlGenerator.GenerateURL(longURL)
		_, err = lr.storage.Get(shortURL)
		attemps += 1
	}

	if attemps == lr.maxGenAttemps && err == nil {
		return "", fmt.Errorf("internal error")
	}

	return lr.storage.Add(shortURL, longURL)
}

func (lr *BasicLinkRepository) ReStore(shortURL string) (string, error) {
	if !lr.urlGenerator.IsValidURL(shortURL) {
		return "", fmt.Errorf("invalid shortURL")
	}
	return lr.storage.Get(shortURL)
}
