package memento

// Originator is an interface that defines the methods required for creating and restoring Mementos.
type Originator interface {
	// CreateMemento creates a new Memento that stores the current state of the object.
	CreateMemento() (*Memento, error)
	// RestoreMemento restores the object's state from a given Memento.
	RestoreMemento(m *Memento) error
}
