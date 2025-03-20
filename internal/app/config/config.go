package config

import (
	"time"
)

// Default app config constants.
const (
	defaultURLGenTimeout       = 2 * time.Second
	defaultURLGenRetryInterval = 100 * time.Millisecond
	serverShutdownTimeout      = 5 * time.Second
	defaultReadHeaderTimeout   = 10 * time.Second  // Maximum duration to read request headers
	defaultWriteTimeout        = 10 * time.Second  // Maximum duration to write response
	defaultIdleTimeout         = 120 * time.Second // Maximum duration for idle connections
	defaultURLSize             = 8
)

// Config holds the app configuration settings, which can be set through environment variables or flags.
//
//easyjson:json
type Config struct {
	ServerAddr              string `env:"SERVER_ADDRESS" json:"server_address"`
	BaseURL                 string `env:"BASE_URL" json:"base_url"`
	FileStoragePath         string `env:"FILE_STORAGE_PATH" json:"file_storage_path"`
	DatabaseDSN             string `env:"DATABASE_DSN" json:"database_dsn"`
	EnableHTTPS             bool   `env:"ENABLE_HTTPS" json:"enable_https"`
	JWTSecret               string `env:"JWT_SECRET" json:"jwt_secret"`
	TLSKeyPath              string `env:"TLC_KEY_PATH" json:"tlc_key_path"`
	TLSCertPath             string `env:"TLC_CERT_PATH" json:"tlc_cert_path"`
	ConfigJSON              string `env:"CONFIG"`
	URLGenTimeout           time.Duration
	URLGenRetryInterval     time.Duration
	URLsize                 int
	ServerShutTimeout       time.Duration
	ServerReadHeaderTimeout time.Duration
	ServerWriteTimeout      time.Duration
	ServerIdleTimeout       time.Duration
	ForceEmptyRepo          bool
}

// DefaultConfig app config.
func DefaultConfig() *Config {
	return &Config{
		ServerAddr:              `localhost:8080`,
		BaseURL:                 `http://localhost:8080/`,
		FileStoragePath:         `data/service_storage.json`,
		DatabaseDSN:             ``,
		EnableHTTPS:             false,
		JWTSecret:               `d1a58c288a0226998149277b14993f6c73cf44ff9df3de548df4df25a13b251a`,
		TLSKeyPath:              `/etc/ssl/private/shortener-key.pem`,
		TLSCertPath:             `/etc/ssl/certs/shortener-cert.pem`,
		ConfigJSON:              ``,
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

// LoadConfig initializes and returns a Config instance.
func LoadConfig() *Config {
	builder := newBuilder()
	builder.loadFile()
	builder.loadEnv()
	builder.loadFlags()
	cfg := builder.getConfig()

	return cfg
}
