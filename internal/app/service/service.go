package service

import (
	"errors"

	"github.com/cenkalti/backoff/v4"
	"github.com/patraden/ya-practicum-go-shortly/internal/app/config"
	e "github.com/patraden/ya-practicum-go-shortly/internal/app/errors"
	"github.com/patraden/ya-practicum-go-shortly/internal/app/repository"
	"github.com/patraden/ya-practicum-go-shortly/internal/app/urlgenerator"
	"github.com/patraden/ya-practicum-go-shortly/internal/app/utils"
)

type URLShortener struct {
	repo         repository.URLRepository
	urlGenerator urlgenerator.URLGenerator
	config       *config.Config
}

func NewShortenerService(
	repo repository.URLRepository,
	gen urlgenerator.URLGenerator,
	config *config.Config,
) *URLShortener {
	return &URLShortener{
		repo:         repo,
		urlGenerator: gen,
		config:       config,
	}
}

func (s *URLShortener) ShortenURL(longURL string) (string, error) {
	// always assume that url generation is an non-injective function.
	// timeout based backoff is the basic mechanism to address collisions.
	// in case of high rates of collisions errors,
	// the intention should rather be to improve URLGenerator algorithms or service.
	var shortURL string
	var err error

	b := utils.LinearBackoff(s.config.URLGenTimeout, s.config.URLGenRetryInterval)
	operation := func() error {
		shortURL = s.urlGenerator.GenerateURL(longURL)
		_, err = s.repo.GetURL(shortURL)

		switch {
		case errors.Is(err, e.ErrRepoNotFound):
			return nil
		case err != nil:
			return backoff.Permanent(err)
		default:
			return e.ErrServiceCollision
		}
	}

	err = backoff.Retry(operation, b)
	if errors.Is(err, e.ErrServiceCollision) {
		return "", e.ErrServiceCollision
	}

	if err != nil {
		return "", e.ErrServiceInternal
	}

	if err = s.repo.AddURL(shortURL, longURL); err != nil {
		return "", e.ErrServiceInternal
	}

	return shortURL, nil
}

func (s *URLShortener) GetOriginalURL(shortURL string) (string, error) {
	if !s.urlGenerator.IsValidURL(shortURL) {
		return "", e.ErrServiceInvalid
	}

	longURL, err := s.repo.GetURL(shortURL)

	if errors.Is(err, e.ErrRepoNotFound) {
		return "", e.ErrRepoNotFound
	}

	if err != nil {
		return "", e.ErrServiceInternal
	}

	return longURL, nil
}
