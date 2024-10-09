package service

import (
	"errors"

	e "github.com/patraden/ya-practicum-go-shortly/internal/app/errors"
	"github.com/patraden/ya-practicum-go-shortly/internal/app/repository"
	"github.com/patraden/ya-practicum-go-shortly/internal/app/urlgenerator"
)

type ShortenerService struct {
	URLShortener
	repo         repository.URLRepository
	urlGenerator urlgenerator.URLGenerator
	retries      int
}

func NewShortenerService() *ShortenerService {
	return &ShortenerService{
		repo:         repository.NewMapURLRepository(),
		urlGenerator: urlgenerator.NewRandURLGenerator(8),
		retries:      1000,
	}
}

func (s *ShortenerService) ShortenURL(longURL string) (string, error) {
	shortURL := s.urlGenerator.GenerateURL(longURL)

	tries := 0
	_, err := s.repo.AddURL(shortURL, longURL)
	for errors.Is(err, e.ErrExists) && tries < s.retries {
		shortURL = s.urlGenerator.GenerateURL(longURL)
		_, err = s.repo.AddURL(shortURL, longURL)
		tries += 1
	}

	if err != nil {
		return "", e.ErrInternal
	}

	return shortURL, nil
}

func (s *ShortenerService) GetOriginalURL(shortURL string) (string, error) {
	if !s.urlGenerator.IsValidURL(shortURL) {
		return "", e.ErrInvalid
	}

	longURL, err := s.repo.GetURL(shortURL)

	if errors.Is(err, e.ErrNotFound) {
		return "", e.ErrNotFound
	}

	if err != nil {
		return "", e.ErrInternal
	}

	return longURL, nil
}
