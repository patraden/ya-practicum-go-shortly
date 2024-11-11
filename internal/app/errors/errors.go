package errors

import (
	"errors"
	"fmt"
)

const (
	WrapOpenFile      = "unable to open file %s: %w"
	WrapCloseFile     = "unable to close file %s: %w"
	WrapUnmarchalJSON = "unable to unmarchal from json: %w"
	WrapMarchalJSON   = "unable to marchal to json: %w"
	WrapFileRead      = "file reader error: %w"
	WrapFileReset     = "file reader reset: %w"
	WrapFileWrite     = "file writer error: %w"
	WrapURLGen        = "URL generate error: %w"
)

var (
	ErrRepoNotFound     = errors.New("repository: URL not found")
	ErrRepoExists       = errors.New("repository: URL exists")
	ErrConfParams       = errors.New("config: invalid config parameter")
	ErrConfEnv          = errors.New("config: env parsing error")
	ErrServiceInternal  = errors.New("service: internal error")
	ErrServiceInvalid   = errors.New("service: invalid URL")
	ErrServiceCollision = errors.New("service: URL collision")
	ErrDecompress       = errors.New("request decompression error")
	ErrTest             = errors.New("testing error")
	ErrUtils            = errors.New("utils error")
	ErrDBOpen           = errors.New("utils db: failed to open db")
	ErrDBPing           = errors.New("utils db: failed to ping db")
	ErrDBDSNParse       = errors.New("utils db: failed to parse dsn")
	ErrDBEmptyPool      = errors.New("databse conn pool is empty")
	ErrDBClose          = errors.New("utils db: close error")
	ErrMementoCreate    = errors.New("memento: could not create astate")
	ErrMementoRestore   = errors.New("memento: could not restore state")
	ErrURLGenerateSlugs = errors.New("urlgenerator: failed to generate slugs")
)

func Wrap(msg string, err error) error {
	return fmt.Errorf("%s: %w", msg, err)
}
