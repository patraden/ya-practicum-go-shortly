package service

import (
	"errors"

	"github.com/cenkalti/backoff/v4"
	"github.com/patraden/ya-practicum-go-shortly/internal/app/config"
	e "github.com/patraden/ya-practicum-go-shortly/internal/app/errors"
	"github.com/patraden/ya-practicum-go-shortly/internal/app/repository"
	"github.com/patraden/ya-practicum-go-shortly/internal/app/repository/file"
	"github.com/patraden/ya-practicum-go-shortly/internal/app/urlgenerator"
	"github.com/patraden/ya-practicum-go-shortly/internal/app/utils"
	"github.com/rs/zerolog"
)

type ShortenerService struct {
	repo         repository.URLRepository
	urlGenerator urlgenerator.URLGenerator
	config       config.Config
}

func NewShortenerService(
	repo repository.URLRepository,
	urlgen urlgenerator.URLGenerator,
	config config.Config,
) *ShortenerService {
	return &ShortenerService{
		repo:         repo,
		urlGenerator: urlgen,
		config:       config,
	}
}

func NewInFileShortenerService(config config.Config, log zerolog.Logger) *ShortenerService {
	return NewShortenerService(
		file.NewInFileURLRepository(config.FileStoragePath, log),
		urlgenerator.NewRandURLGenerator(config.URLsize),
		config,
	)
}

func NewInMemoryShortenerService(config config.Config) *ShortenerService {
	return NewShortenerService(
		repository.NewInMemoryURLRepository(),
		urlgenerator.NewRandURLGenerator(config.URLsize),
		config,
	)
}

func (s *ShortenerService) ShortenURL(longURL string) (string, error) {
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
		case errors.Is(err, e.ErrNotFound):
			return nil
		case err != nil:
			return backoff.Permanent(err)
		default:
			return e.ErrCollision
		}
	}

	err = backoff.Retry(operation, b)
	if errors.Is(err, e.ErrCollision) {
		return "", e.ErrCollision
	}

	if err != nil {
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
