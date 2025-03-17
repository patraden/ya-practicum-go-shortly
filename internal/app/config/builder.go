package config

import (
	"flag"
	"log"
	"strings"

	"github.com/caarlos0/env/v6"

	e "github.com/patraden/ya-practicum-go-shortly/internal/app/domain/errors"
	"github.com/patraden/ya-practicum-go-shortly/internal/app/utils"
)

type builder struct {
	cfg *Config
}

func newBuilder() *builder {
	return &builder{
		cfg: DefaultConfig(),
	}
}

func (b *builder) loadEnv() {
	if err := env.Parse(b.cfg); err != nil {
		log.Fatal(e.ErrEnvConfigParse)
	}
}

func (b *builder) loadFlags() {
	flag.StringVar(&b.cfg.ServerAddr, "a", b.cfg.ServerAddr, "server address {host}:{port}")
	flag.StringVar(&b.cfg.BaseURL, "b", b.cfg.BaseURL, "base url {base url}/{short link}")
	flag.StringVar(&b.cfg.FileStoragePath, "f", b.cfg.FileStoragePath, "url storage file path")
	flag.StringVar(&b.cfg.DatabaseDSN, "d", b.cfg.DatabaseDSN, "database DSN")
	flag.BoolVar(&b.cfg.EnableHTTPS, "s", b.cfg.EnableHTTPS, "enable https")
	flag.BoolVar(&b.cfg.ForceEmptyRepo, "force-empty", b.cfg.ForceEmptyRepo, "skip repository load from disk")
	flag.Parse()
}

func (b *builder) getConfig() *Config {
	if !utils.IsServerAddress(b.cfg.ServerAddr) {
		log.Fatal(e.ErrInvalidConfig)
	}

	if !strings.HasSuffix(b.cfg.BaseURL, "/") {
		b.cfg.BaseURL += "/"
	}

	return b.cfg
}
