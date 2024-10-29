package config

import (
	"flag"
	"log"
	"strings"

	"github.com/caarlos0/env/v6"
	e "github.com/patraden/ya-practicum-go-shortly/internal/app/errors"
	"github.com/patraden/ya-practicum-go-shortly/internal/app/utils"
)

type builder struct {
	env   *Config
	flags *Config
}

func newBuilder() *builder {
	return &builder{
		env:   DefaultConfig(),
		flags: DefaultConfig(),
	}
}

func (b *builder) loadEnvConfig() {
	err := env.Parse(b.env)
	if err != nil {
		log.Fatal(e.ErrConfEnv)
	}
}

func (b *builder) loadFlagsConfig() {
	cfg := DefaultConfig()

	flag.StringVar(&b.flags.ServerAddr, "a", cfg.ServerAddr, "server address {host}:{port}")
	flag.StringVar(&b.flags.BaseURL, "b", cfg.BaseURL, "base url {base url}/{short link}")
	flag.StringVar(&b.flags.FileStoragePath, "f", cfg.FileStoragePath, "url storage file path")
	flag.Parse()
}

func (b *builder) getConfig() *Config {
	cfg := DefaultConfig()

	// handle Server Address
	switch {
	case b.env.ServerAddr != cfg.ServerAddr:
		cfg.ServerAddr = b.env.ServerAddr
	case b.flags.ServerAddr != cfg.ServerAddr:
		cfg.ServerAddr = b.flags.ServerAddr
	}

	// validate Server Address
	if !utils.IsServerAddress(cfg.ServerAddr) {
		log.Fatal(e.ErrConfParams)
	}

	// handle Base URL
	switch {
	case b.env.BaseURL != cfg.BaseURL:
		cfg.BaseURL = b.env.BaseURL
	case b.flags.BaseURL != cfg.BaseURL:
		cfg.BaseURL = b.flags.BaseURL
	}

	// validate Base URL
	if !utils.IsURL(cfg.BaseURL) {
		log.Fatal(e.ErrConfParams)
	}

	if !strings.HasSuffix(cfg.BaseURL, "/") {
		cfg.BaseURL += "/"
	}

	// handle Base URL
	switch {
	case b.env.FileStoragePath != cfg.FileStoragePath:
		cfg.FileStoragePath = b.env.FileStoragePath
	case b.flags.FileStoragePath != cfg.FileStoragePath:
		cfg.FileStoragePath = b.flags.FileStoragePath
	}

	return cfg
}
