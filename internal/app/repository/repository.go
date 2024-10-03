package repository

import (
	"errors"
	"fmt"
)

// Initial idea is to keep storage and repository as two separate services.
// Mainly because repository will be composed of different components
// like variations of URL generator service and storages
type LinkRepository interface {
	Store(longURL string) (string, error)
	ReStore(shortURL string) (string, error)
}

const erroFormat = "repository error: %w"

var (
	ErrInternal = fmt.Errorf(erroFormat, errors.New("internal error"))
	ErrNotFound = fmt.Errorf(erroFormat, errors.New("url not found"))
	ErrInvalid  = fmt.Errorf(erroFormat, errors.New("invalid url format"))
)
