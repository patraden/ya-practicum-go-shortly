package config_test

import (
	"flag"
	"os"
	"testing"

	easyjson "github.com/mailru/easyjson"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/patraden/ya-practicum-go-shortly/internal/app/config"
)

func resetFlags(t *testing.T) {
	t.Helper()

	flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ContinueOnError)
}

func TestDefaultConfig(t *testing.T) {
	t.Parallel()

	cfg := config.DefaultConfig()

	if cfg.ServerAddr != "localhost:8080" {
		t.Errorf("Expected default ServerAddr 'localhost:8080', got '%s'", cfg.ServerAddr)
	}

	if cfg.EnableHTTPS {
		t.Errorf("Expected default EnableHTTPS to be false")
	}
}

func TestLoadConfigFromFile(t *testing.T) {
	t.Parallel()

	tmpFile, err := os.CreateTemp(t.TempDir(), "config_test.json")
	if err != nil {
		t.Fatal(err)
	}

	defer os.Remove(tmpFile.Name())

	jsonData := `{
		"server_address": "192.168.1.100:9000",
		"enable_https": true
	}`

	if _, err := tmpFile.WriteString(jsonData); err != nil {
		t.Fatal(err)
	}

	tmpFile.Close()

	os.Args = []string{"shortener", "-config", tmpFile.Name()}

	resetFlags(t)

	cfg := config.LoadConfig()
	if cfg.ServerAddr != "192.168.1.100:9000" {
		t.Errorf("Expected ServerAddr '192.168.1.100:9000', got '%s'", cfg.ServerAddr)
	}

	if !cfg.EnableHTTPS {
		t.Errorf("Expected EnableHTTPS to be true, got false")
	}
}

// Test priority order: Env > Flags > File > Default.
func TestLoadConfigPriorityOrder(t *testing.T) {
	tmpFile, err := os.CreateTemp(t.TempDir(), "config_test.json")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpFile.Name())

	jsonData := `{
		"server_address": "json-config:9090",
		"enable_https": false
	}`

	if _, err := tmpFile.WriteString(jsonData); err != nil {
		t.Fatal(err)
	}

	tmpFile.Close()

	t.Setenv("SERVER_ADDRESS", "env-config:8080")
	t.Setenv("ENABLE_HTTPS", "true")

	resetFlags(t)

	os.Args = []string{"cmd", "-c", tmpFile.Name(), "-a", "flag-config:7000", "-s"}

	flag.Parse()

	cfg := config.LoadConfig()

	// Env hould take priority
	if cfg.ServerAddr != "env-config:8080" {
		t.Errorf("Expected ServerAddr 'env-config:8080', got '%s'", cfg.ServerAddr)
	}

	if !cfg.EnableHTTPS {
		t.Errorf("Expected EnableHTTPS to be true, got false")
	}
}

func TestLoadConfigSaveAsFile(t *testing.T) {
	t.Parallel()

	originalCfg := config.DefaultConfig()

	jsonData, err := easyjson.Marshal(originalCfg)
	require.NoError(t, err, "failed to marshal config")

	tmpFile, err := os.CreateTemp(t.TempDir(), "config-*.json")
	require.NoError(t, err, "failed to create temp file")

	defer os.Remove(tmpFile.Name())

	_, err = tmpFile.Write(jsonData)
	require.NoError(t, err, "failed to write config to file")

	tmpFile.Close()

	data, err := os.ReadFile(tmpFile.Name())
	require.NoError(t, err, "failed to read config file")

	var loadedCfg config.Config
	err = easyjson.Unmarshal(data, &loadedCfg)
	require.NoError(t, err, "failed to unmarshal config")

	assert.Equal(t, *originalCfg, loadedCfg, "config mismatch after unmarshaling")
}
