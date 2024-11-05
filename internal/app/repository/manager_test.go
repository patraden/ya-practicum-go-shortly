package repository_test

import (
	"os"
	"testing"

	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/patraden/ya-practicum-go-shortly/internal/app/config"
	"github.com/patraden/ya-practicum-go-shortly/internal/app/logger"
	"github.com/patraden/ya-practicum-go-shortly/internal/app/repository"
)

func TestFileReadWrite(t *testing.T) {
	t.Parallel()

	testRepoState := repository.NewURLRepositoryState(
		map[string]string{
			"XXXYZZZZ": "http://ya.com",
			"XXXYYZZZ": "http://ya.com",
			"XXXYYYZZ": "http://ya.com",
		},
	)

	t.Run("test save/load service state", func(t *testing.T) {
		t.Parallel()

		log := logger.NewLogger(zerolog.InfoLevel).GetLogger()
		cfg := config.DefaultConfig()

		cfg.FileStoragePath = "records.json"
		defer os.Remove(cfg.FileStoragePath)

		manager := repository.NewStateManager(cfg, log)

		err := manager.SaveToFile(testRepoState)
		require.NoError(t, err)

		repoState, err := manager.LoadFromFile()
		require.NoError(t, err)
		assert.Equal(t, len(testRepoState.GetState()), len(repoState.GetState()))
	})
}
