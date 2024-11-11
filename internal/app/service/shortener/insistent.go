package shortener

import (
	"context"
	"errors"
	"time"

	"github.com/cenkalti/backoff/v4"
	"github.com/rs/zerolog"

	"github.com/patraden/ya-practicum-go-shortly/internal/app/config"
	"github.com/patraden/ya-practicum-go-shortly/internal/app/domain"
	"github.com/patraden/ya-practicum-go-shortly/internal/app/dto"
	e "github.com/patraden/ya-practicum-go-shortly/internal/app/errors"
	"github.com/patraden/ya-practicum-go-shortly/internal/app/repository"
	"github.com/patraden/ya-practicum-go-shortly/internal/app/service/urlgenerator"
	"github.com/patraden/ya-practicum-go-shortly/internal/app/utils"
)

const batchGenFactor = 100

type InsistentShortener struct {
	repo         repository.URLRepository
	urlGenerator urlgenerator.URLGenerator
	config       *config.Config
	log          zerolog.Logger
}

func NewInsistentShortener(
	repo repository.URLRepository,
	gen urlgenerator.URLGenerator,
	config *config.Config,
	log zerolog.Logger,
) *InsistentShortener {
	return &InsistentShortener{
		repo:         repo,
		urlGenerator: gen,
		config:       config,
		log:          log,
	}
}

func (s *InsistentShortener) backoff(ctx context.Context) *backoff.ExponentialBackOff {
	// always assume that url generation is an non-injective function.
	// timeout based backoff is the basic mechanism to address collisions.
	// in case of high rates of collisions errors,
	// the intention should rather be to improve URLGenerator algorithms or service.
	b := utils.LinearBackoff(s.config.URLGenTimeout, s.config.URLGenRetryInterval)
	backoff.WithContext(b, ctx)

	return b
}

func (s *InsistentShortener) ShortenURL(ctx context.Context, original domain.OriginalURL) (*domain.URLMapping, error) {
	var slug domain.Slug
	var err error

	boff := s.backoff(ctx)
	try := 1
	operation := func() error {
		s.log.
			Info().
			Int("try", try).
			Msg("trying to shorten URL batch")

		try++

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

	err = backoff.Retry(operation, boff)
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

func (s *InsistentShortener) ShortenURLBatch(ctx context.Context, batch *dto.OriginalURLBatch) (dto.SlugBatch, error) {
	size := len(*batch)
	originals := batch.Originals()
	urlMappings, res := make([]domain.URLMapping, size), make(dto.SlugBatch, size)

	boff, try := s.backoff(ctx), 1
	operation := func() error {
		s.log.Info().Int("try", try).Msg("trying to shorten URL batch")

		try++

		ctxWithTO, cancel := context.WithTimeout(ctx, time.Duration(batchGenFactor*size)*time.Millisecond)
		defer cancel()

		slugs, err := s.urlGenerator.GenerateSlugs(ctxWithTO, originals)
		if errors.Is(err, e.ErrURLGenerateSlugs) {
			return backoff.Permanent(err)
		}

		for i, slug := range slugs {
			urlMappings[i] = *domain.NewURLMapping(slug, originals[i])
		}

		if err = s.repo.AddURLMappingBatch(ctx, &urlMappings); err == nil {
			if len(slugs) != size {
				return e.ErrURLGenerateSlugs
			}

			for i, elem := range *batch {
				res[i] = dto.CorrelatedSlug{CorrelationID: elem.CorrelationID, Slug: slugs[i].WithBaseURL(s.config.BaseURL)}
			}

			return nil
		}

		if errors.Is(err, e.ErrRepoExists) {
			return e.ErrServiceCollision
		}

		return e.Wrap("db repo batch:", err)
	}

	if err := backoff.Retry(operation, boff); err != nil {
		s.log.Error().Err(err).Msg("shorten batch failed")

		return dto.SlugBatch{}, e.ErrServiceInternal
	}

	return res, nil
}
