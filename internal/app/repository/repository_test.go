package repository_test

import (
	"errors"
	"os"
	"testing"

	e "github.com/patraden/ya-practicum-go-shortly/internal/app/errors"
	"github.com/patraden/ya-practicum-go-shortly/internal/app/logger"
	"github.com/patraden/ya-practicum-go-shortly/internal/app/repository"
	"github.com/patraden/ya-practicum-go-shortly/internal/app/repository/file"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/require"
)

type want struct {
	err     error
	longURL string
}
type test struct {
	name     string
	longURL  string
	shortURL string
	want     want
}

func RunRepositoryTests(t *testing.T, repo repository.URLRepository) {
	t.Helper()

	addURLTests := []test{
		{
			name:     "test 1",
			shortURL: "a",
			longURL:  "b",
			want: want{
				err:     nil,
				longURL: "",
			},
		},
		{
			name:     "test 2",
			shortURL: "b",
			longURL:  "c",
			want: want{
				err:     nil,
				longURL: "",
			},
		},
		{
			name:     "test 3",
			shortURL: "a",
			longURL:  "b",
			want: want{
				err:     e.ErrExists,
				longURL: "",
			},
		},
	}

	for _, test := range addURLTests {
		t.Run(test.name, func(t *testing.T) {
			err := repo.AddURL(test.shortURL, test.longURL)
			if !errors.Is(err, test.want.err) {
				t.Errorf("AddURL test failed: shortURL %s, longURL %s, err: %v", test.shortURL, test.longURL, err)
			}
		})
	}

	getURLTests := []test{
		{
			name:     "test 1",
			shortURL: "a",
			longURL:  "b",
			want: want{
				err:     nil,
				longURL: "b",
			},
		},
		{
			name:     "test 2",
			shortURL: "b",
			longURL:  "c",
			want: want{
				err:     nil,
				longURL: "c",
			},
		},
		{
			name:     "test 3",
			shortURL: "c",
			longURL:  "d",
			want: want{
				err:     e.ErrNotFound,
				longURL: "",
			},
		},
	}

	for _, test := range getURLTests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			longURL, err := repo.GetURL(test.shortURL)
			if !errors.Is(err, test.want.err) || longURL != test.want.longURL {
				t.Errorf("GetURL test failed: shortURL %s, longURL %s, err: %v", test.shortURL, test.longURL, err)
			}
		})
	}
}

func TestInMemoryURLRepository(t *testing.T) {
	t.Parallel()

	repo := repository.NewInMemoryURLRepository()
	RunRepositoryTests(t, repo)
}

func TestInFileURLRepository(t *testing.T) {
	t.Parallel()

	fileName := "file/records.json"
	defer os.Remove(fileName)

	log := logger.NewLogger(zerolog.InfoLevel).GetLogger()
	repo := file.NewInFileURLRepository(fileName, log)
	require.NotNil(t, repo)
	RunRepositoryTests(t, repo)
}
