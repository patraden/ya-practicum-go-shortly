package config

import (
	"flag"
	"log"
	"os"
	"strings"

	"github.com/caarlos0/env/v6"
	easyjson "github.com/mailru/easyjson"

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

func (b *builder) loadFile() {
	args := os.Args[1:]

	for i := range args {
		if args[i] == "-c" || args[i] == "-config" {
			if i+1 < len(args) {
				b.cfg.ConfigJSON = args[i+1]
			}

			break
		}
	}

	if b.cfg.ConfigJSON == `` {
		return
	}

	file, err := os.ReadFile(b.cfg.ConfigJSON)
	if err != nil {
		log.Fatal(e.ErrInvalidConfig)

		return
	}

	if err := easyjson.Unmarshal(file, b.cfg); err != nil {
		log.Fatal(e.ErrInvalidConfig)

		return
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
	flag.StringVar(&b.cfg.TrustedSubnet, "t", b.cfg.TrustedSubnet, "trusted subnet")
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
