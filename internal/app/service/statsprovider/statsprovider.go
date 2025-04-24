package statsprovider

import (
	"context"

	"github.com/patraden/ya-practicum-go-shortly/internal/app/dto"
)

// StatsProvider is an interface that collects variaous repo statistics.
type StatsProvider interface {
	GetStats(ctx context.Context) (*dto.RepoStats, error)
}
