package repository

import (
	"errors"

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
	shortURL := lr.urlGenerator.GenerateURL(longURL)

	attemps := 0
	_, err := lr.storage.Add(shortURL, longURL)
	for err != nil && errors.Is(err, storage.ErrKeyExists) && attemps < lr.maxGenAttemps {
		shortURL = lr.urlGenerator.GenerateURL(longURL)
		_, err = lr.storage.Add(shortURL, longURL)
		attemps += 1
	}

	if err != nil {
		if errors.Is(err, storage.ErrKeyExists) {
			return "", ErrOutOfURL
		}
		return "", ErrInternal
	}
	return shortURL, nil
}

func (lr *BasicLinkRepository) ReStore(shortURL string) (string, error) {
	if !lr.urlGenerator.IsValidURL(shortURL) {
		return "", ErrInvalid
	}

	longURL, err := lr.storage.Get(shortURL)

	if err != nil {
		if errors.Is(err, storage.ErrKeyNotFound) {
			return "", ErrNotFound
		}
		return "", ErrInternal
	}
	return longURL, nil
}
