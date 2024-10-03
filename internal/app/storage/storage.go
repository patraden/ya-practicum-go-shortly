package storage

import (
	"errors"
	"fmt"
)

// Basic Key Value storage interface
type BasicKVStorage interface {
	Add(key string, value string) (string, error)
	Get(key string) (string, error)
	Del(key string) error
}

const erroFormat = "storage error: %w"

var (
	ErrKeyNotFound = fmt.Errorf(erroFormat, errors.New("key not found"))
	ErrKeyExists   = fmt.Errorf(erroFormat, errors.New("key already exists"))
)
