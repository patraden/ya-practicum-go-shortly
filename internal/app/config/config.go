package config

import "time"

type Config struct {
	ServerAddr    string        `env:"SERVER_ADDRESS"`
	BaseURL       string        `env:"BASE_URL"`
	URLGenTimeout time.Duration `env:"URL_GENENERATE_TIMEOUT"`
}

func DefaultConfig() Config {
	return Config{
		ServerAddr:    `localhost:8080`,
		BaseURL:       `http://localhost:8080/`,
		URLGenTimeout: time.Duration(10) * time.Second,
	}
}

func LoadConfig() Config {
	builder := newBuilder()
	builder.loadEnvConfig()
	builder.loadFlagsConfig()
	cfg := builder.getConfig()
	return cfg
}
