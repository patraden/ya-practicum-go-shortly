package memento

import (
	"bufio"
	"io"
	"os"

	"github.com/mailru/easyjson"
	"github.com/rs/zerolog"

	"github.com/patraden/ya-practicum-go-shortly/internal/app/config"
	"github.com/patraden/ya-practicum-go-shortly/internal/app/domain"
	e "github.com/patraden/ya-practicum-go-shortly/internal/app/domain/errors"
	"github.com/patraden/ya-practicum-go-shortly/internal/app/dto"
)

const (
	PermReadWriteUser = 0o644 // Read/write for owner, read-only for others
	EOL               = "\n"
	errLabel          = "memento"
)

type StateManager struct {
	config     *config.Config
	originator Originator
	log        *zerolog.Logger
}

func NewStateManager(config *config.Config, originator Originator, log *zerolog.Logger) *StateManager {
	return &StateManager{
		config:     config,
		originator: originator,
		log:        log,
	}
}

func (sm *StateManager) RestoreFromState(state *Memento) error {
	if err := sm.originator.RestoreMemento(state); err != nil {
		return e.ErrStateRestore
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
		return e.ErrStateRestore
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
		return e.ErrStateCreate
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
	log     *zerolog.Logger
}

func NewReader(fileName string, log *zerolog.Logger) (*Reader, error) {
	file, err := os.OpenFile(fileName, os.O_RDONLY|os.O_CREATE, PermReadWriteUser)
	if err != nil {
		log.Error().
			Err(err).
			Str("filename", fileName).
			Msg("failed to open file")

		return nil, e.Wrap("failed to open file", err, errLabel)
	}

	return &Reader{
		file:    file,
		scanner: bufio.NewScanner(file),
		log:     log,
	}, nil
}

func (r *Reader) LoadState() (*Memento, error) {
	state := make(dto.URLMappings)
	var count int

	if err := r.Reset(); err != nil {
		return nil, err
	}

	for r.scanner.Scan() {
		data := r.scanner.Bytes()
		link := domain.URLMapping{}

		err := link.UnmarshalJSON(data)
		if err != nil {
			r.log.Error().
				Err(err).
				Msg("failed to unmarchal state")

			return nil, e.Wrap("failed to unmarshal state", err, errLabel)
		}

		state[link.Slug] = link
		count++

		r.log.Info().
			Str("short_url", string(link.Slug)).
			Str("long_url", string(link.OriginalURL)).
			Str("user_id", link.UserID.String()).
			Time("created_at", link.CreatedAt).
			Time("expires_at", link.ExpiresAt).
			Msg("loaded record")
	}

	if err := r.scanner.Err(); err != nil {
		return nil, e.Wrap("file scanner error", err, errLabel)
	}

	r.log.Info().
		Int("total_records", count).
		Msg("completed loading state")

	return NewMemento(state), nil
}

func (r *Reader) Reset() error {
	if _, err := r.file.Seek(0, io.SeekStart); err != nil {
		return e.Wrap("file scanner seek error", err, errLabel)
	}

	r.scanner = bufio.NewScanner(r.file)

	return nil
}

func (r *Reader) Close() error {
	err := r.file.Close()
	if err != nil {
		return e.Wrap("failed to close file", err, errLabel)
	}

	return nil
}

type Writer struct {
	file *os.File
	log  *zerolog.Logger
}

func NewWriter(fileName string, log *zerolog.Logger) (*Writer, error) {
	file, err := os.OpenFile(fileName, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, PermReadWriteUser)
	if err != nil {
		log.Error().
			Err(err).
			Str("filename", fileName).
			Msg("failed to open file")

		return nil, e.Wrap("failed to open file", err, errLabel)
	}

	return &Writer{
		file: file,
		log:  log,
	}, nil
}

func (w *Writer) SaveState(state *Memento) error {
	writer := bufio.NewWriter(w.file)
	defer writer.Flush()

	var count int

	for _, link := range state.GetState() {
		if _, err := easyjson.MarshalToWriter(link, writer); err != nil {
			w.log.
				Error().
				Err(err).
				Msg("failed to write state")

			return e.Wrap("failed to write state", err, errLabel)
		}

		if _, err := writer.WriteString(EOL); err != nil {
			return e.Wrap("failed to write EOL", err, errLabel)
		}

		count++

		w.log.Info().
			Str("short_url", string(link.Slug)).
			Str("long_url", string(link.OriginalURL)).
			Str("user_id", link.UserID.String()).
			Time("created_at", link.CreatedAt).
			Time("expires_at", link.ExpiresAt).
			Msg("preserved record")
	}

	w.log.Info().
		Int("total_records", count).
		Msg("completed saving state")

	return nil
}

func (w *Writer) Close() error {
	err := w.file.Close()
	if err != nil {
		return e.Wrap("failed to close file", err, errLabel)
	}

	return nil
}
