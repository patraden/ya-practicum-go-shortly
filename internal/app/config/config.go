package config

import "time"

const (
	defaultURLGenTimeout       = time.Duration(2) * time.Second
	defaultURLGenRetryInterval = time.Duration(100) * time.Millisecond
	defaultURLSize             = 8
)

type Config struct {
	ServerAddr          string `env:"SERVER_ADDRESS"`
	BaseURL             string `env:"BASE_URL"`
	FileStoragePath     string `env:"FILE_STORAGE_PATH"`
	URLGenTimeout       time.Duration
	URLGenRetryInterval time.Duration
	URLsize             int
}

func DefaultConfig() Config {
	return Config{
		ServerAddr:          `localhost:8080`,
		BaseURL:             `http://localhost:8080/`,
		FileStoragePath:     `data/records.json`,
		URLGenTimeout:       defaultURLGenTimeout,
		URLGenRetryInterval: defaultURLGenRetryInterval,
		URLsize:             defaultURLSize,
	}
}

func LoadConfig() Config {
	builder := newBuilder()
	builder.loadEnvConfig()
	builder.loadFlagsConfig()
	cfg := builder.getConfig()

	return cfg
}
