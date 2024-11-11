package memento

type Originator interface {
	CreateMemento() (*Memento, error)
	RestoreMemento(m *Memento) error
}
