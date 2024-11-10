package shortener

import (
	"context"
	"errors"

	"github.com/cenkalti/backoff/v4"

	"github.com/patraden/ya-practicum-go-shortly/internal/app/config"
	"github.com/patraden/ya-practicum-go-shortly/internal/app/domain"
	e "github.com/patraden/ya-practicum-go-shortly/internal/app/errors"
	"github.com/patraden/ya-practicum-go-shortly/internal/app/repository"
	"github.com/patraden/ya-practicum-go-shortly/internal/app/service/urlgenerator"
	"github.com/patraden/ya-practicum-go-shortly/internal/app/utils"
)

type InsistentShortener struct {
	repo         repository.URLRepository
	urlGenerator urlgenerator.URLGenerator
	config       *config.Config
}

func NewInsistentShortener(
	repo repository.URLRepository,
	gen urlgenerator.URLGenerator,
	config *config.Config,
) *InsistentShortener {
	return &InsistentShortener{
		repo:         repo,
		urlGenerator: gen,
		config:       config,
	}
}

func (s *InsistentShortener) ShortenURL(ctx context.Context, original domain.OriginalURL) (*domain.URLMapping, error) {
	// always assume that url generation is an non-injective function.
	// timeout based backoff is the basic mechanism to address collisions.
	// in case of high rates of collisions errors,
	// the intention should rather be to improve URLGenerator algorithms or service.
	var slug domain.Slug
	var err error

	b := utils.LinearBackoff(s.config.URLGenTimeout, s.config.URLGenRetryInterval)
	operation := func() error {
		slug = s.urlGenerator.GenerateSlug(ctx, original)
		_, err = s.repo.GetURLMapping(ctx, slug)

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
		return nil, e.ErrServiceCollision
	}

	if err != nil {
		return nil, e.ErrServiceInternal
	}

	m := domain.NewURLMapping(slug, original)

	if err = s.repo.AddURLMapping(ctx, m); err != nil {
		return nil, e.ErrServiceInternal
	}

	return m, nil
}

func (s *InsistentShortener) GetOriginalURL(ctx context.Context, slug domain.Slug) (*domain.URLMapping, error) {
	if !s.urlGenerator.IsValidSlug(slug) {
		return nil, e.ErrServiceInvalid
	}

	m, err := s.repo.GetURLMapping(ctx, slug)

	if errors.Is(err, e.ErrRepoNotFound) {
		return nil, e.ErrRepoNotFound
	}

	if err != nil {
		return nil, e.ErrServiceInternal
	}

	return m, nil
}
