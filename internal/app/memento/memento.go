package memento

import "github.com/patraden/ya-practicum-go-shortly/internal/app/dto"

type Memento struct {
	state dto.URLMappings
}

func NewMemento(state dto.URLMappings) *Memento {
	return &Memento{
		state: state,
	}
}

func (m *Memento) GetState() dto.URLMappings {
	return m.state
}
