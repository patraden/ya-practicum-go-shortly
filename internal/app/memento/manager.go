package memento

import (
	"bufio"
	"fmt"
	"io"
	"os"

	"github.com/mailru/easyjson"
	"github.com/rs/zerolog"

	"github.com/patraden/ya-practicum-go-shortly/internal/app/config"
	"github.com/patraden/ya-practicum-go-shortly/internal/app/domain"
	"github.com/patraden/ya-practicum-go-shortly/internal/app/dto"
	e "github.com/patraden/ya-practicum-go-shortly/internal/app/errors"
)

const (
	PermReadWriteUser = 0o644 // Read/write for owner, read-only for others
	EOL               = "\n"
)

type StateManager struct {
	config     *config.Config
	originator Originator
	log        zerolog.Logger
}

func NewStateManager(config *config.Config, originator Originator, log zerolog.Logger) *StateManager {
	return &StateManager{
		config:     config,
		originator: originator,
		log:        log,
	}
}

func (sm *StateManager) RestoreFromState(state *Memento) error {
	if err := sm.originator.RestoreMemento(state); err != nil {
		return e.ErrMementoRestore
	}

	return nil
}

func (sm *StateManager) RestoreFromFile() error {
	r, err := NewReader(sm.config.FileStoragePath, sm.log)
	if err != nil {
		return err
	}
	defer r.Close()

	state, err := r.LoadState()
	if err != nil {
		return err
	}

	if err := sm.originator.RestoreMemento(state); err != nil {
		return e.ErrMementoRestore
	}

	return nil
}

func (sm *StateManager) StoreToFile() error {
	w, err := NewWriter(sm.config.FileStoragePath, sm.log)
	if err != nil {
		return err
	}
	defer w.Close()

	state, err := sm.originator.CreateMemento()
	if err != nil {
		return e.ErrMementoCreate
	}

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
	state := make(dto.URLMappings)

	if err := r.Reset(); err != nil {
		return nil, err
	}

	for r.scanner.Scan() {
		data := r.scanner.Bytes()
		link := domain.URLMapping{}

		err := link.UnmarshalJSON(data)
		if err != nil {
			return nil, fmt.Errorf(e.WrapUnmarchalJSON, err)
		}

		state[link.Slug] = link
		r.log.Info().
			Str("short_url", string(link.Slug)).
			Str("long_url", string(link.OriginalURL)).
			Time("created_at", link.CreatedAt).
			Time("expires_at", link.ExpiresAt).
			Msg("loaded record")
	}

	if r.scanner.Err() == nil {
		return NewMemento(state), nil
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
	for _, link := range state.GetState() {
		if _, err := easyjson.MarshalToWriter(link, w.file); err != nil {
			return fmt.Errorf(e.WrapMarchalJSON, err)
		}

		if _, err := w.file.WriteString(EOL); err != nil {
			return fmt.Errorf(e.WrapFileWrite, err)
		}

		w.log.Info().
			Str("short_url", string(link.Slug)).
			Str("long_url", string(link.OriginalURL)).
			Time("created_at", link.CreatedAt).
			Time("expires_at", link.ExpiresAt).
			Msg("preserved record")
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
