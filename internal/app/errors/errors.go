package errors

import (
	"errors"
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
)

type GeneralError struct {
	Err error
}

func (e *GeneralError) Error() string {
	return e.Err.Error()
}

func (e *GeneralError) Unwrap() error {
	return e.Err
}

func (e *GeneralError) Is(target error) bool {
	_, ok := target.(*GeneralError)

	return ok
}

func (e *GeneralError) As(target interface{}) bool {
	if t, ok := target.(*GeneralError); ok {
		*t = *e

		return true
	}

	return false
}

func General(err error) error {
	if err == nil {
		return nil
	}

	return &GeneralError{
		Err: err,
	}
}
