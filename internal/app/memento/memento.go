package memento

import "github.com/patraden/ya-practicum-go-shortly/internal/app/dto"

// Memento is a struct that stores a snapshot of the URL mappings state.
type Memento struct {
	state dto.URLMappings
}

// NewMemento creates and returns a new Memento instance with the given state.
func NewMemento(state dto.URLMappings) *Memento {
	return &Memento{
		state: state,
	}
}

// GetState returns the current stored state of the Memento.
func (m *Memento) GetState() dto.URLMappings {
	return m.state
}
