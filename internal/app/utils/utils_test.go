package utils_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/patraden/ya-practicum-go-shortly/internal/app/utils"
)

const alphabet = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

func TestIsServerAddress(t *testing.T) {
	t.Parallel()

	validAddr := "localhost:8080"
	assert.True(t, utils.IsServerAddress(validAddr))

	invalidAddr := "localhost"
	assert.False(t, utils.IsServerAddress(invalidAddr))
}

func TestLinearBackoff(t *testing.T) {
	t.Parallel()

	maxElapsedTime := 10 * time.Second
	interval := 1 * time.Second

	bo := utils.LinearBackoff(maxElapsedTime, interval)

	assert.NotNil(t, bo)
	assert.Equal(t, maxElapsedTime, bo.MaxElapsedTime)
	assert.Equal(t, interval, bo.InitialInterval)
}

func TestRandString(t *testing.T) {
	t.Parallel()

	length := 16
	randomString := utils.RandString(length)
	assert.Len(t, randomString, length)

	for _, char := range randomString {
		assert.Contains(t, alphabet, string(char))
	}
}

func TestRandInt(t *testing.T) {
	t.Parallel()

	mx := 100
	randomInt := utils.RandInt(mx)
	assert.True(t, randomInt >= 0 && randomInt < mx)
}

func TestRandURL(t *testing.T) {
	t.Parallel()

	randomURL := utils.RandURL()

	assert.Contains(t, randomURL, "://")
	assert.Contains(t, randomURL, ".")
	assert.Contains(t, randomURL, "resource/")
	assert.Contains(t, randomURL, "id=")
}
