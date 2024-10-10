package config

import (
	"flag"
	"log"
	"strings"

	"github.com/caarlos0/env/v6"
	e "github.com/patraden/ya-practicum-go-shortly/internal/app/errors"
	"github.com/patraden/ya-practicum-go-shortly/internal/app/helpers"
)

type defaultBuilder struct {
	env   Config
	flags Config
}

func newDefaultBuilder() *defaultBuilder {
	return &defaultBuilder{}
}

func (b *defaultBuilder) loadEnvConfig() {
	err := env.Parse(&b.env)
	if err != nil {
		log.Fatal(e.ErrEnvParse)
	}

}

func (b *defaultBuilder) loadFlagsConfig() {
	flag.StringVar(&b.flags.ServerAddr, "a", "", "server address {host}:{port}")
	flag.StringVar(&b.flags.BaseURL, "b", "", "base url {base url}/{short link}")
	flag.Parse()
}

func (b *defaultBuilder) getConfig() Config {
	cfg := DefaultConfig()

	// handle Server Address
	switch {
	case b.env.ServerAddr != "":
		cfg.ServerAddr = b.env.ServerAddr
	case b.flags.ServerAddr != "":
		cfg.ServerAddr = b.flags.ServerAddr
	}

	// validate Server Address
	if !helpers.IsServerAddress(cfg.ServerAddr) {
		log.Fatal(e.ErrParams)
	}

	// handle Base URL
	switch {
	case b.env.BaseURL != "":
		cfg.BaseURL = b.env.BaseURL
	case b.flags.BaseURL != "":
		cfg.BaseURL = b.flags.BaseURL
	}

	// validate Base URL
	if !helpers.IsURL(cfg.BaseURL) {
		log.Fatal(e.ErrParams)
	}

	if !strings.HasSuffix(cfg.BaseURL, "/") {
		cfg.BaseURL += "/"
	}

	return cfg

}
