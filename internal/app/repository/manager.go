package repository

import (
	"bufio"
	"fmt"
	"io"
	"os"

	easyjson "github.com/mailru/easyjson"
	"github.com/patraden/ya-practicum-go-shortly/internal/app/config"
	e "github.com/patraden/ya-practicum-go-shortly/internal/app/errors"
	"github.com/rs/zerolog"
)

const (
	PermReadWriteUser = 0o644 // Read/write for owner, read-only for others
	EOL               = "\n"
)

type StateManager struct {
	config *config.Config
	log    zerolog.Logger
}

func NewStateManager(config *config.Config, log zerolog.Logger) *StateManager {
	return &StateManager{
		log:    log,
		config: config,
	}
}

func (sm *StateManager) LoadFromFile() (*Memento, error) {
	r, err := NewReader(sm.config.FileStoragePath, sm.log)
	if err != nil {
		return nil, err
	}
	defer r.Close()

	state, err := r.LoadState()
	if err != nil {
		return nil, err
	}

	return state, nil
}

func (sm *StateManager) SaveToFile(state *Memento) error {
	w, err := NewWriter(sm.config.FileStoragePath, sm.log)
	if err != nil {
		return err
	}
	defer w.Close()

	err = w.SaveState(state)
	if err != nil {
		return err
	}

	return nil
}

type Reader struct {
	file    *os.File
	scanner *bufio.Scanner
	log     zerolog.Logger
}

func NewReader(fileName string, log zerolog.Logger) (*Reader, error) {
	file, err := os.OpenFile(fileName, os.O_RDONLY|os.O_CREATE, PermReadWriteUser)
	if err != nil {
		return nil, fmt.Errorf(e.WrapOpenFile, fileName, err)
	}

	return &Reader{
		file:    file,
		scanner: bufio.NewScanner(file),
		log:     log,
	}, nil
}

func (r *Reader) LoadState() (*Memento, error) {
	state := make(map[string]string)

	if err := r.Reset(); err != nil {
		return nil, err
	}

	for r.scanner.Scan() {
		data := r.scanner.Bytes()
		record := FileRecord{}

		err := record.UnmarshalJSON(data)
		if err != nil {
			return nil, fmt.Errorf(e.WrapUnmarchalJSON, err)
		}

		state[record.ShortURL] = record.LongURL
		r.log.Info().
			Int("uuid", record.ID).
			Str("short_url", record.ShortURL).
			Str("long_url", record.LongURL).
			Msg("loaded record")
	}

	if r.scanner.Err() == nil {
		return NewURLRepositoryState(state), nil
	}

	return nil, fmt.Errorf(e.WrapFileRead, r.scanner.Err())
}

func (r *Reader) Reset() error {
	if _, err := r.file.Seek(0, io.SeekStart); err != nil {
		return fmt.Errorf(e.WrapFileReset, r.scanner.Err())
	}

	r.scanner = bufio.NewScanner(r.file)

	return nil
}

func (r *Reader) Close() error {
	err := r.file.Close()
	if err != nil {
		return fmt.Errorf(e.WrapCloseFile, r.file.Name(), err)
	}

	return nil
}

type Writer struct {
	file *os.File
	log  zerolog.Logger
}

func NewWriter(fileName string, log zerolog.Logger) (*Writer, error) {
	file, err := os.OpenFile(fileName, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, PermReadWriteUser)
	if err != nil {
		return nil, fmt.Errorf(e.WrapOpenFile, fileName, err)
	}

	return &Writer{
		file: file,
		log:  log,
	}, nil
}

func (w *Writer) SaveState(state *Memento) error {
	urls := state.GetState()

	index := 1
	for shortURL, longURL := range urls {
		record := &FileRecord{
			ID:       index,
			ShortURL: shortURL,
			LongURL:  longURL,
		}

		if _, err := easyjson.MarshalToWriter(record, w.file); err != nil {
			return fmt.Errorf(e.WrapMarchalJSON, err)
		}

		if _, err := w.file.WriteString(EOL); err != nil {
			return fmt.Errorf(e.WrapFileWrite, err)
		}

		w.log.Info().
			Int("uuid", record.ID).
			Str("short_url", record.ShortURL).
			Str("long_url", record.LongURL).
			Msg("preserved record")

		index++
	}

	return nil
}

func (w *Writer) Close() error {
	err := w.file.Close()
	if err != nil {
		return fmt.Errorf(e.WrapCloseFile, w.file.Name(), err)
	}

	return nil
}
