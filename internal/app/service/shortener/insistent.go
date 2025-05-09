package shortener

import (
	"context"
	"errors"
	"time"

	"github.com/cenkalti/backoff/v4"
	"github.com/rs/zerolog"

	"github.com/patraden/ya-practicum-go-shortly/internal/app/config"
	"github.com/patraden/ya-practicum-go-shortly/internal/app/domain"
	e "github.com/patraden/ya-practicum-go-shortly/internal/app/domain/errors"
	"github.com/patraden/ya-practicum-go-shortly/internal/app/dto"
	"github.com/patraden/ya-practicum-go-shortly/internal/app/middleware"
	"github.com/patraden/ya-practicum-go-shortly/internal/app/repository"
	"github.com/patraden/ya-practicum-go-shortly/internal/app/service/urlgenerator"
	"github.com/patraden/ya-practicum-go-shortly/internal/app/utils"
)

const batchGenFactor = 100

// InsistentShortener is a URL shortener that retries generating slugs in case of collisions
// with backoff logic to ensure successful URL shortening.
type InsistentShortener struct {
	repo         repository.URLRepository
	urlGenerator urlgenerator.URLGenerator
	config       *config.Config
	log          *zerolog.Logger
}

// NewInsistentShortener creates a new instance of InsistentShortener with the provided repository,
// URL generator, configuration, and logger.
func NewInsistentShortener(
	repo repository.URLRepository,
	gen urlgenerator.URLGenerator,
	config *config.Config,
	log *zerolog.Logger,
) *InsistentShortener {
	return &InsistentShortener{
		repo:         repo,
		urlGenerator: gen,
		config:       config,
		log:          log,
	}
}

func (s *InsistentShortener) generateSlugWithBackoff(ctx context.Context, operation func() error) error {
	// always assume that url generation is an non-injective function.
	// timeout based backoff is the basic mechanism to address collisions.
	// in case of high rates of collisions errors,
	// the intention should rather be to improve URLGenerator algorithms or service.
	boff := utils.LinearBackoff(s.config.URLGenTimeout, s.config.URLGenRetryInterval)
	backoff.WithContext(boff, ctx)

	try := 1

	err := backoff.Retry(func() error {
		try++

		return operation()
	}, boff)
	if err != nil {
		return e.Wrap("retry error", err, errLabel)
	}

	return nil
}

// ShortenURL shortens a URL by generating a unique slug for the original URL and saving it to the repository.
// It retries generating a slug in case of collisions.
func (s *InsistentShortener) ShortenURL(ctx context.Context, original domain.OriginalURL) (domain.Slug, error) {
	var slug domain.Slug

	operation := func() error {
		slug = s.urlGenerator.GenerateSlug(ctx, original)
		_, errRepo := s.repo.GetURLMapping(ctx, slug)

		if errors.Is(errRepo, e.ErrSlugNotFound) {
			return nil
		}

		if errRepo != nil {
			return backoff.Permanent(errRepo)
		}

		return e.ErrSlugCollision
	}

	if err := s.generateSlugWithBackoff(ctx, operation); err != nil {
		if errors.Is(err, e.ErrSlugCollision) {
			return "", e.ErrSlugCollision
		}

		s.log.Error().Err(err).Msg("slug generation failed")

		return "", e.ErrShortenerInternal
	}

	userID, ok := middleware.GetUserID(ctx)
	if !ok {
		s.log.Error().Msg("failed to get userID from context")

		return "", e.ErrShortenerInternal
	}

	newMap := domain.NewURLMapping(slug, original, userID)
	m, err := s.repo.AddURLMapping(ctx, newMap)

	if errors.Is(err, e.ErrOriginalExists) {
		return m.Slug, e.ErrOriginalExists
	}

	if err != nil {
		s.log.Error().Err(err).Msg("failed to shorten url")

		return "", e.ErrShortenerInternal
	}

	return m.Slug, nil
}

// GetOriginalURL retrieves the original URL associated with the given slug.
// If the slug does not exist or has been deleted, appropriate errors are returned.
func (s *InsistentShortener) GetOriginalURL(ctx context.Context, slug domain.Slug) (domain.OriginalURL, error) {
	if !s.urlGenerator.IsValidSlug(slug) {
		return "", e.ErrSlugInvalid
	}

	urlm, err := s.repo.GetURLMapping(ctx, slug)

	if errors.Is(err, e.ErrSlugNotFound) {
		return "", e.ErrSlugNotFound
	}

	if err != nil {
		s.log.Error().Err(err).Msg("shortener internal error")

		return "", e.ErrShortenerInternal
	}

	if urlm.Deleted {
		return urlm.OriginalURL, e.ErrSlugDeleted
	}

	return urlm.OriginalURL, nil
}

// GetUserURLs retrieves all URL mappings for a specific user.
// It returns the URLs in a batch format for efficiency.
func (s *InsistentShortener) GetUserURLs(ctx context.Context) (*dto.URLPairBatch, error) {
	userID, ok := middleware.GetUserID(ctx)
	if !ok {
		s.log.Error().Msg("failed to get userID from context")

		return &dto.URLPairBatch{}, e.ErrShortenerInternal
	}

	res, err := s.repo.GetUserURLMappings(ctx, userID)

	if errors.Is(err, e.ErrUserNotFound) {
		return &dto.URLPairBatch{}, e.ErrUserNotFound
	}

	if err != nil {
		s.log.Error().Err(err).Msg("shortener internal error")

		return &dto.URLPairBatch{}, e.ErrShortenerInternal
	}

	return dto.NewURLPairBatch(&res, s.config.BaseURL), nil
}

// ShortenURLBatch shortens a batch of URLs by generating unique slugs for each one and storing the mappings.
// It retries generating slugs in case of collisions for the batch of URLs.
func (s *InsistentShortener) ShortenURLBatch(ctx context.Context, batch *dto.OriginalURLBatch) (*dto.SlugBatch, error) {
	size := len(*batch)
	originals := batch.Originals()
	urlMappings := make([]domain.URLMapping, size)
	res := make(dto.SlugBatch, size)

	userID, ok := middleware.GetUserID(ctx)
	if !ok {
		s.log.Error().Msg("failed to get userID from context")

		return &dto.SlugBatch{}, e.ErrShortenerInternal
	}

	operation := func() error {
		ctxWithTO, cancel := context.WithTimeout(ctx, time.Duration(batchGenFactor*size)*time.Millisecond)
		defer cancel()
		// generating slugs for batch
		slugs, err := s.urlGenerator.GenerateSlugs(ctxWithTO, originals)
		if errors.Is(err, e.ErrURLGenGenerateSlug) {
			// cannot generate unique set of slugs - stop retrying
			return backoff.Permanent(err)
		}
		// generating a batch of urlmappings
		for i, slug := range slugs {
			urlMappings[i] = *domain.NewURLMapping(slug, originals[i], userID)
		}
		// trying to add them to repo
		err = s.repo.AddURLMappingBatch(ctx, &urlMappings)
		if err == nil {
			if len(slugs) != size {
				return e.ErrURLGenGenerateSlug
			}

			for i, elem := range *batch {
				res[i] = dto.CorrelatedSlug{CorrelationID: elem.CorrelationID, Slug: slugs[i]}
			}
			// success - stop retrying
			return nil
		}
		// collisions - continue retrying
		if errors.Is(err, e.ErrSlugExists) {
			return e.ErrSlugCollision
		}
		// unexpected error - stop retrying
		return backoff.Permanent(err)
	}

	if err := s.generateSlugWithBackoff(ctx, operation); err != nil {
		if errors.Is(err, e.ErrSlugCollision) {
			return nil, e.ErrSlugCollision
		}

		s.log.Error().Err(err).Msg("failed to shorten url batch")

		return &dto.SlugBatch{}, e.ErrShortenerInternal
	}

	return &res, nil
}
