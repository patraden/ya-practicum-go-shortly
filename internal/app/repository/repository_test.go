package repository_test

import (
	"errors"
	"testing"

	e "github.com/patraden/ya-practicum-go-shortly/internal/app/errors"
	"github.com/patraden/ya-practicum-go-shortly/internal/app/repository"
	"github.com/stretchr/testify/assert"
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
				err:     e.ErrRepoExists,
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
				err:     e.ErrRepoNotFound,
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

	t.Run("memento test", func(t *testing.T) {
		t.Parallel()

		memento, err := repo.CreateMemento()
		require.NoError(t, err)

		state := memento.GetState()
		before := len(state)

		state["newURL"] = "newURL"

		memento = repository.NewURLRepositoryState(state)
		err = repo.RestoreMemento(memento)
		require.NoError(t, err)

		memento, err = repo.CreateMemento()
		require.NoError(t, err)

		state = memento.GetState()
		after := len(state)

		assert.Equal(t, before+1, after)
	})
}

func TestInMemoryURLRepository(t *testing.T) {
	t.Parallel()

	repo := repository.NewInMemoryURLRepository()
	RunRepositoryTests(t, repo)
}
