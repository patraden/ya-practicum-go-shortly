package file

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"sync"

	easyjson "github.com/mailru/easyjson"
	e "github.com/patraden/ya-practicum-go-shortly/internal/app/errors"
	"github.com/patraden/ya-practicum-go-shortly/internal/app/utils"
	"github.com/rs/zerolog"
)

const (
	PermReadWriteUser = 0o644 // Read/write for owner, read-only for others
	EOL               = "\n"
)

// To avoid offset conflicts
// two file descriptors will be maintained in Reader and Writer.
// Besides, repository will include naive map based inmemory cache.
type InFileURLRepository struct {
	sync.RWMutex
	cache  map[string]string
	reader *Reader
	writer *Writer
	log    zerolog.Logger
}

func NewInFileURLRepository(fileName string, log zerolog.Logger) *InFileURLRepository {
	reader := NewReader(fileName, log)
	if reader == nil {
		return nil
	}

	writer := NewWriter(fileName, log)
	if writer == nil {
		return nil
	}

	cache := make(map[string]string)
	if err := reader.LoadToCache(&cache); err != nil {
		return nil
	}

	return &InFileURLRepository{
		RWMutex: sync.RWMutex{},
		cache:   cache,
		reader:  reader,
		writer:  writer,
		log:     log,
	}
}

func (fs *InFileURLRepository) AddURL(shortURL string, longURL string) error {
	// serializable transactions isolation :)
	fs.Lock()
	defer fs.Unlock()

	if _, ok := fs.cache[shortURL]; ok {
		return e.ErrExists
	}

	if record, err := fs.reader.Find(shortURL); err == nil {
		// add missing shortURL to cache
		fs.cache[shortURL] = record.LongURL

		return e.ErrExists
	}

	rec := &Record{
		UUID:     utils.UUID(),
		ShortURL: shortURL,
		LongURL:  longURL,
	}

	// store into file first
	if err := fs.writer.Write(rec); err != nil {
		return e.ErrRepoFile
	}

	fs.cache[shortURL] = longURL

	return nil
}

func (fs *InFileURLRepository) GetURL(shortURL string) (string, error) {
	fs.Lock()
	defer fs.Unlock()

	value, ok := fs.cache[shortURL]
	if ok {
		return value, nil
	}

	if record, err := fs.reader.Find(shortURL); err == nil {
		// add missing shortURL to cache
		fs.cache[shortURL] = record.LongURL

		return record.LongURL, nil
	}

	return "", e.ErrNotFound
}

func (fs *InFileURLRepository) Close() error {
	if err := fs.reader.Close(); err != nil {
		return err
	}

	if err := fs.writer.Close(); err != nil {
		return err
	}

	return nil
}

func (fs *InFileURLRepository) DelURL(_ string) error {
	return nil
}

type Reader struct {
	file    *os.File
	scanner *bufio.Scanner
	log     zerolog.Logger
}

func NewReader(fileName string, log zerolog.Logger) *Reader {
	file, err := os.OpenFile(fileName, os.O_RDONLY|os.O_CREATE, PermReadWriteUser)
	if err != nil {
		return nil
	}

	return &Reader{
		file:    file,
		scanner: bufio.NewScanner(file),
		log:     log,
	}
}

func (r *Reader) Read() (*Record, error) {
	if !r.scanner.Scan() {
		if r.scanner.Err() == nil {
			return nil, io.EOF
		}

		return nil, e.Wrap(r.scanner.Err(), e.ErrRepoFile)
	}

	data := r.scanner.Bytes()
	rec := Record{}

	err := rec.UnmarshalJSON(data)
	if err != nil {
		return nil, e.Wrap(err, e.ErrRepoFile)
	}

	return &rec, nil
}

// terrible full scan find for now.
func (r *Reader) Find(shortURL string) (*Record, error) {
	if err := r.Reset(); err != nil {
		return nil, err
	}

	for {
		record, err := r.Read()
		if err != nil {
			if errors.Is(err, io.EOF) {
				return nil, fmt.Errorf("%w: record not found", e.ErrRepoFile)
			}

			return nil, err
		}

		if record.ShortURL == shortURL {
			return record, nil
		}
	}
}

func (r *Reader) LoadToCache(cache *map[string]string) error {
	if err := r.Reset(); err != nil {
		return err
	}

	for r.scanner.Scan() {
		data := r.scanner.Bytes()
		record := Record{}

		err := record.UnmarshalJSON(data)
		if err != nil {
			return e.Wrap(err, e.ErrRepoFile)
		}

		_, ok := (*cache)[record.ShortURL]
		if ok {
			r.log.Info().
				Str("uuid", record.UUID).
				Str("short_url", record.ShortURL).
				Str("long_url", record.LongURL).
				Msg("duplicate record")

			continue
		}

		(*cache)[record.ShortURL] = record.LongURL
		r.log.Info().
			Str("uuid", record.UUID).
			Str("short_url", record.ShortURL).
			Str("long_url", record.LongURL).
			Msg("loaded to cache")
	}

	if r.scanner.Err() == nil {
		return nil
	}

	return e.Wrap(r.scanner.Err(), e.ErrRepoFile)
}

func (r *Reader) Reset() error {
	if _, err := r.file.Seek(0, io.SeekStart); err != nil {
		return e.Wrap(err, e.ErrRepoFile)
	}

	r.scanner = bufio.NewScanner(r.file)

	return nil
}

func (r *Reader) Close() error {
	err := r.file.Close()
	if err != nil {
		return e.Wrap(err, e.ErrRepoFile)
	}

	return nil
}

type Writer struct {
	file *os.File
	log  zerolog.Logger
}

func NewWriter(fileName string, log zerolog.Logger) *Writer {
	file, err := os.OpenFile(fileName, os.O_WRONLY|os.O_CREATE|os.O_APPEND, PermReadWriteUser)
	if err != nil {
		return nil
	}

	return &Writer{
		file: file,
		log:  log,
	}
}

func (w *Writer) Write(record *Record) error {
	if _, err := easyjson.MarshalToWriter(record, w.file); err != nil {
		return e.Wrap(err, e.ErrRepoFile)
	}

	if _, err := w.file.WriteString(EOL); err != nil {
		return e.Wrap(err, e.ErrRepoFile)
	}

	w.log.Info().
		Str("uuid", record.UUID).
		Str("short_url", record.ShortURL).
		Str("long_url", record.LongURL).
		Msg("preserved record")

	return nil
}

func (w *Writer) Close() error {
	err := w.file.Close()
	if err != nil {
		return e.Wrap(err, e.ErrRepoFile)
	}

	return nil
}
