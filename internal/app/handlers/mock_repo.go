package handlers

import (
	"github.com/patraden/ya-practicum-go-shortly/internal/app/repository"
	"github.com/stretchr/testify/mock"
)

type MockLinkRepository struct {
	mock.Mock
	repository.LinkRepository
}

func (m *MockLinkRepository) Store(longURL string) (string, error) {
	args := m.Called(longURL)
	return args.String(0), args.Error(1)
}

func (m *MockLinkRepository) ReStore(shortURL string) (string, error) {
	args := m.Called(shortURL)
	return args.String(0), args.Error(1)
}
