package errors

import (
	"errors"
	"fmt"
)

var (
	ErrParams     = errors.New("invalid config parameter")
	ErrEnvParse   = errors.New("env config parsing error")
	ErrNotFound   = errors.New("URL not found in repository")
	ErrExists     = errors.New("URL exists in repository")
	ErrInternal   = errors.New("internal URL service error")
	ErrInvalid    = errors.New("invalid URL format")
	ErrDecompress = errors.New("request decompression error")
	ErrCollision  = errors.New("URL collision")
	ErrTest       = errors.New("testing error")
	ErrUtils      = errors.New("utils error")
	ErrRepoFile   = errors.New("file repo error")
)

func Wrap(isErr error, asErr error) error {
	return fmt.Errorf("%s: %w", asErr.Error(), isErr)
}
