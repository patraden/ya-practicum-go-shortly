package repository

type Memento struct {
	state map[string]string
}

func NewURLRepositoryState(state map[string]string) *Memento {
	return &Memento{
		state: state,
	}
}

func (m *Memento) GetState() map[string]string {
	return m.state
}
