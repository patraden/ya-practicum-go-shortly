package memento_test

import (
	"os"
	"testing"

	"github.com/rs/zerolog"
	"github.com/stretchr/testify/require"

	"github.com/patraden/ya-practicum-go-shortly/internal/app/config"
	"github.com/patraden/ya-practicum-go-shortly/internal/app/domain"
	"github.com/patraden/ya-practicum-go-shortly/internal/app/dto"
	"github.com/patraden/ya-practicum-go-shortly/internal/app/logger"
	"github.com/patraden/ya-practicum-go-shortly/internal/app/memento"
	"github.com/patraden/ya-practicum-go-shortly/internal/app/repository"
)

func TestFileReadWrite(t *testing.T) {
	t.Parallel()

	originator := repository.NewInMemoryURLRepository()

	userID := domain.NewUserID()
	testMemento := memento.NewMemento(
		dto.URLMappings{
			"XXXYZZZZ": *domain.NewURLMapping("XXXYZZZZ", "http://ya.com", userID),
			"XXXYYZZZ": *domain.NewURLMapping("XXXYYZZZ", "http://ya.com", domain.NewUserID()),
			"XXXYYYZZ": *domain.NewURLMapping("XXXYYYZZ", "http://ya.com", domain.NewUserID()),
		},
	)

	t.Run("test save/load service state", func(t *testing.T) {
		t.Parallel()

		log := logger.NewLogger(zerolog.InfoLevel).GetLogger()
		cfg := config.DefaultConfig()

		cfg.FileStoragePath = "records.json"
		defer os.Remove(cfg.FileStoragePath)

		manager := memento.NewStateManager(cfg, originator, log)
		err := manager.RestoreFromState(testMemento)
		require.NoError(t, err)

		err = manager.StoreToFile()
		require.NoError(t, err)

		err = manager.RestoreFromFile()
		require.NoError(t, err)
	})
}
