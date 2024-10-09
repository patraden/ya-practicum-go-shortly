package errors

import (
	"errors"
	"fmt"
)

var (
	ErrParams   = fmt.Errorf(`config error: %w`, errors.New("invalid parameter"))
	ErrEnvParse = fmt.Errorf(`config error: %w`, errors.New("env parsing error"))

	ErrNotFound = fmt.Errorf(`repository error: %w`, errors.New("url not found"))
	ErrExists   = fmt.Errorf(`repository error: %w`, errors.New("url already exist"))

	ErrInternal = fmt.Errorf(`service error: %w`, errors.New("internal error"))
	ErrInvalid  = fmt.Errorf(`service error: %w`, errors.New("invalid url format"))
)
