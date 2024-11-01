package config

import (
	"flag"
	"log"
	"reflect"
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
	if err := env.Parse(b.env); err != nil {
		log.Fatal(e.ErrConfEnv)
	}
}

func (b *builder) loadFlagsConfig() {
	flag.StringVar(&b.flags.ServerAddr, "a", b.flags.ServerAddr, "server address {host}:{port}")
	flag.StringVar(&b.flags.BaseURL, "b", b.flags.BaseURL, "base url {base url}/{short link}")
	flag.StringVar(&b.flags.FileStoragePath, "f", b.flags.FileStoragePath, "url storage file path")
	flag.BoolVar(&b.flags.ForceEmptyRepo, "force-empty", false, "do not load and preserve repository")
	flag.Parse()
}

func (b *builder) getConfig() *Config {
	cfg := DefaultConfig()

	applyField := func(field string) {
		envValue := reflect.ValueOf(b.env).Elem().FieldByName(field)
		flagValue := reflect.ValueOf(b.flags).Elem().FieldByName(field)
		cfgValue := reflect.ValueOf(cfg).Elem().FieldByName(field)

		// Prioritize environment variable, then flags, then default.
		switch {
		case envValue.String() != cfgValue.String():
			cfgValue.Set(envValue)
		case flagValue.String() != cfgValue.String():
			cfgValue.Set(flagValue)
		}
	}

	fields := []string{"ServerAddr", "BaseURL", "FileStoragePath"}
	for _, field := range fields {
		applyField(field)
	}

	cfg.ForceEmptyRepo = b.flags.ForceEmptyRepo

	if !utils.IsServerAddress(cfg.ServerAddr) {
		log.Fatal(e.ErrConfParams)
	}

	if !utils.IsURL(cfg.BaseURL) {
		log.Fatal(e.ErrConfParams)
	}

	if !strings.HasSuffix(cfg.BaseURL, "/") {
		cfg.BaseURL += "/"
	}

	return cfg
}
