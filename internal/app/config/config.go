package config

import (
	"time"
)

const (
	defaultURLGenTimeout       = 2 * time.Second
	defaultURLGenRetryInterval = 100 * time.Millisecond
	serverShutdownTimeout      = 5 * time.Second
	defaultReadHeaderTimeout   = 10 * time.Second  // Maximum duration to read request headers
	defaultWriteTimeout        = 10 * time.Second  // Maximum duration to write response
	defaultIdleTimeout         = 120 * time.Second // Maximum duration for idle connections
	defaultURLSize             = 8
)

type Config struct {
	ServerAddr              string `env:"SERVER_ADDRESS"`
	BaseURL                 string `env:"BASE_URL"`
	FileStoragePath         string `env:"FILE_STORAGE_PATH"`
	DatabaseDSN             string `env:"DATABASE_DSN"`
	URLGenTimeout           time.Duration
	URLGenRetryInterval     time.Duration
	URLsize                 int
	ServerShutTimeout       time.Duration
	ServerReadHeaderTimeout time.Duration
	ServerWriteTimeout      time.Duration
	ServerIdleTimeout       time.Duration
	ForceEmptyRepo          bool
}

func DefaultConfig() *Config {
	return &Config{
		ServerAddr:              `localhost:8080`,
		BaseURL:                 `http://localhost:8080/`,
		FileStoragePath:         `data/service_storage.json`,
		DatabaseDSN:             ``,
		URLGenTimeout:           defaultURLGenTimeout,
		URLGenRetryInterval:     defaultURLGenRetryInterval,
		URLsize:                 defaultURLSize,
		ServerShutTimeout:       serverShutdownTimeout,
		ServerReadHeaderTimeout: defaultReadHeaderTimeout,
		ServerWriteTimeout:      defaultWriteTimeout,
		ServerIdleTimeout:       defaultIdleTimeout,
		ForceEmptyRepo:          false,
	}
}

func LoadConfig() *Config {
	builder := newBuilder()
	builder.loadEnvConfig()
	builder.loadFlagsConfig()
	cfg := builder.getConfig()

	return cfg
}
