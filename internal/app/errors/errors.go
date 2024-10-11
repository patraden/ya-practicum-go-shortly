package errors

import (
	"errors"
)

var (
	ErrParams   = errors.New("invalid config parameter")
	ErrEnvParse = errors.New("env config parsing error")
	ErrNotFound = errors.New("URL not found in repository")
	ErrExists   = errors.New("URL exists in repository")
	ErrInternal = errors.New("internal URL service error")
	ErrInvalid  = errors.New("invalid URL format")
)
