package service

import (
	"errors"
	"time"

	"github.com/cenkalti/backoff/v4"
	e "github.com/patraden/ya-practicum-go-shortly/internal/app/errors"
	"github.com/patraden/ya-practicum-go-shortly/internal/app/repository"
	"github.com/patraden/ya-practicum-go-shortly/internal/app/urlgenerator"
)

const shortURLsize = 8

type ShortenerService struct {
	repo         repository.URLRepository
	urlGenerator urlgenerator.URLGenerator
	genTimeout   time.Duration
}

func NewShortenerService(timeout time.Duration) *ShortenerService {
	return &ShortenerService{
		repo:         repository.NewInMemoryURLRepository(),
		urlGenerator: urlgenerator.NewRandURLGenerator(shortURLsize),
		genTimeout:   timeout,
	}
}

func (s *ShortenerService) ShortenURL(longURL string) (string, error) {
	// always assume that url generation is an non-injective function.
	// timeout based backoff is the basic mechanism to address collisions.
	// in case of high rates of collisions errors,
	// the intention should rather be to improve URLGenerator algorithms
	var shortURL string
	var err error

	op := func() error {
		shortURL = s.urlGenerator.GenerateURL(longURL)
		_, err = s.repo.GetURL(shortURL)

		switch {
		case errors.Is(err, e.ErrNotFound):
			return nil
		case err != nil:
			return backoff.Permanent(err)
		default:
			return e.ErrCollision
		}
	}

	b := backoff.NewExponentialBackOff(backoff.WithMaxElapsedTime(s.genTimeout))
	if err = backoff.Retry(op, b); err != nil {
		return "", e.ErrInternal
	}

	if err = s.repo.AddURL(shortURL, longURL); err != nil {
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
