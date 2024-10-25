package mock

import (
	"github.com/stretchr/testify/mock"
)

type ShortenerService struct {
	mock.Mock
}

func NewShortenerService() *ShortenerService {
	return &ShortenerService{}
}

func (m *ShortenerService) ShortenURL(longURL string) (string, error) {
	args := m.Called(longURL)

	return args.String(0), args.Error(1)
}

func (m *ShortenerService) GetOriginalURL(shortURL string) (string, error) {
	args := m.Called(shortURL)

	return args.String(0), args.Error(1)
}

// use go mock from uber
