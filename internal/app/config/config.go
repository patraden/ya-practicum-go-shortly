package config

import (
	"errors"
	"flag"
	"fmt"
	"log"
	"strings"

	"github.com/caarlos0/env/v6"
	"github.com/patraden/ya-practicum-go-shortly/internal/app/helpers"
	r "github.com/patraden/ya-practicum-go-shortly/internal/app/repository"
)

const (
	defaultServerAddr = `localhost:8080`
	defaultBaseURL    = `http://localhost:8080/`
	errorFormat       = `config error: %w`
)

var (
	ErrServerAddr = fmt.Errorf(errorFormat, errors.New("not a valid server address {host}:{port}"))
	ErrBaseURL    = fmt.Errorf(errorFormat, errors.New("not a valid base url {base url}/{short link}"))
	ErrEnvParse   = fmt.Errorf(errorFormat, errors.New("env parsing error"))
)

type Config struct {
	ServerAddr string `env:"SERVER_ADDRESS"`
	BaseURL    string `env:"BASE_URL"`
	Repo       r.LinkRepository
}

func DefaultConfig(repo r.LinkRepository) *Config {
	return &Config{
		ServerAddr: defaultServerAddr,
		BaseURL:    defaultBaseURL,
		Repo:       repo,
	}
}

func LoadConfig(repo r.LinkRepository) *Config {
	cfg := Config{
		Repo: repo,
	}
	ecfg := Config{}
	fcfg := Config{}

	// load from env
	err := env.Parse(&ecfg)
	if err != nil {
		log.Fatal(ErrEnvParse)
	}

	// load from flags
	flag.StringVar(&fcfg.ServerAddr, "a", defaultServerAddr, "server address {host}:{port}")
	flag.StringVar(&fcfg.BaseURL, "b", defaultBaseURL, "base url {base url}/{short link}")
	flag.Parse()

	// handle Server Address
	switch {
	case ecfg.ServerAddr != "" && !helpers.IsServerAddress(ecfg.ServerAddr):
		log.Fatal(ErrServerAddr)
	case fcfg.ServerAddr != "" && !helpers.IsServerAddress(fcfg.ServerAddr):
		log.Fatal(ErrServerAddr)
	case ecfg.ServerAddr != "":
		cfg.ServerAddr = ecfg.ServerAddr
	case fcfg.ServerAddr != "":
		cfg.ServerAddr = fcfg.ServerAddr
	default:
		cfg.ServerAddr = defaultServerAddr
	}

	// handle Base URL
	switch {
	case ecfg.BaseURL != "" && !helpers.IsURL(ecfg.BaseURL):
		log.Fatal(ErrBaseURL)
	case fcfg.BaseURL != "" && !helpers.IsURL(fcfg.BaseURL):
		log.Fatal(ErrBaseURL)
	case ecfg.BaseURL != "":
		cfg.BaseURL = ecfg.BaseURL
	case fcfg.BaseURL != "":
		cfg.BaseURL = fcfg.BaseURL
	default:
		cfg.BaseURL = defaultBaseURL
	}

	// make sure it always ends with "/"
	if !strings.HasSuffix(cfg.BaseURL, "/") {
		cfg.BaseURL += "/"
	}

	return &cfg
}
