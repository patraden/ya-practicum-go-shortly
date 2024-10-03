package config

import (
	"flag"
	"fmt"
	"strings"

	"github.com/patraden/ya-practicum-go-shortly/internal/app/repository"
)

type Config struct {
	ServerAddr   ServerAddress
	ShortURLAddr string
	Repo         repository.LinkRepository
}

type ServerAddress struct {
	Host string
	Port string
}

func (s ServerAddress) String() string {
	return s.Host + ":" + s.Port
}

func (s *ServerAddress) Set(value string) error {
	hp := strings.Split(value, ":")
	if len(hp) != 2 {
		return fmt.Errorf("need address in a form host:port")
	}
	s.Host = hp[0]
	s.Port = hp[1]
	return nil
}

func NewDevConfig(repo repository.LinkRepository) *Config {
	return &Config{
		ServerAddr: ServerAddress{
			Host: `localhost`,
			Port: `8080`,
		},
		ShortURLAddr: `http://localhost:8000/`,
		Repo:         repo,
	}

}

func DevConfigWithFlags(repo repository.LinkRepository) *Config {
	devConfig := NewDevConfig(repo)
	flag.Var(&devConfig.ServerAddr, "a", "server address {host:port}")
	flag.StringVar(&devConfig.ShortURLAddr, "b", "http://localhost:8000/", "URL address http://localhost:8000/{shortURL}")
	flag.Parse()
	return devConfig
}
