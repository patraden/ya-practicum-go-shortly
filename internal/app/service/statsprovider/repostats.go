package statsprovider

import (
	"context"

	"github.com/rs/zerolog"

	e "github.com/patraden/ya-practicum-go-shortly/internal/app/domain/errors"
	"github.com/patraden/ya-practicum-go-shortly/internal/app/dto"
	"github.com/patraden/ya-practicum-go-shortly/internal/app/repository"
)

// RepoStatsProvider provides statistical data related to URL shortener service.
type RepoStatsProvider struct {
	repo repository.URLRepository
	log  *zerolog.Logger
}

// NewRepoStatsProvider creates a new instance of RepoStatsProvider.
func NewRepoStatsProvider(repo repository.URLRepository, log *zerolog.Logger) *RepoStatsProvider {
	return &RepoStatsProvider{
		repo: repo,
		log:  log,
	}
}

// GetStats retrieves repo statistics.
func (srv *RepoStatsProvider) GetStats(ctx context.Context) (*dto.RepoStats, error) {
	stats, err := srv.repo.GetStats(ctx)
	if err != nil {
		srv.log.Error().Err(err).
			Msg("Failed to get statistics from repo")

		return nil, e.ErrStatsProviderInternal
	}

	return stats, nil
}
