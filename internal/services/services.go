package services

import (
	"fmt"

	"github.com/patraden/ya-practicum-go-shortly/internal/storage"
	"github.com/patraden/ya-practicum-go-shortly/internal/urlgen"
)

type LinkStore struct {
	maxGenAttemps int
	urlGenerator  urlgen.UrlGenerator
	storage       storage.BasicKVStorage
}

func NewSimpleLinkStore() *LinkStore {
	return &LinkStore{
		maxGenAttemps: 1000,
		urlGenerator:  urlgen.NewRandUrlGenerator(8),
		storage:       storage.NewMapStorage(),
	}
}

func (ls *LinkStore) Store(longUrl string) (string, error) {
	attemps := 0
	shortUrl := ls.urlGenerator.GenerateShortUrl(longUrl)
	_, err := ls.storage.Get(shortUrl)

	for err == nil && attemps < ls.maxGenAttemps {
		shortUrl := ls.urlGenerator.GenerateShortUrl(longUrl)
		_, err = ls.storage.Get(shortUrl)
		attemps += 1
	}

	if attemps == ls.maxGenAttemps && err == nil {
		return "", fmt.Errorf("out of shortlinks")
	}

	return ls.storage.Add(shortUrl, longUrl)
}

func (ls *LinkStore) ReStore(shortUrl string) (string, error) {
	return ls.storage.Get(shortUrl)
}
