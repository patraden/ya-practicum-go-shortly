package repository_test

import (
	"errors"
	"testing"

	e "github.com/patraden/ya-practicum-go-shortly/internal/app/errors"
	"github.com/patraden/ya-practicum-go-shortly/internal/app/repository"
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

func RunAddURLRepositoryTests(t *testing.T, repo repository.URLRepository, test test) {
	t.Helper()

	t.Run(test.name, func(t *testing.T) {
		if err := repo.AddURL(test.shortURL, test.longURL); !errors.Is(err, test.want.err) {
			t.Errorf("AddURL test failed: shortURL %s, longURL %s", test.shortURL, test.longURL)
		}
	})
}

func RunGetURLRepositoryTests(t *testing.T, repo repository.URLRepository, test test) {
	t.Helper()

	t.Run(test.name, func(t *testing.T) {
		t.Parallel()

		if longURL, err := repo.GetURL(test.shortURL); !errors.Is(err, test.want.err) || longURL != test.want.longURL {
			t.Errorf("GetURL test failed: shortURL %s, longURL %s", test.shortURL, test.longURL)
		}
	})
}

func TestInMemoryURLRepository(t *testing.T) {
	t.Parallel()

	repo := repository.NewInMemoryURLRepository()

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
		RunAddURLRepositoryTests(t, repo, test)
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
		RunGetURLRepositoryTests(t, repo, test)
	}
}
