package repository_test

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

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

func testAddURL(t *testing.T, repo repository.URLRepository) {
	t.Helper()

	addURLTests := []test{
		{"test 1", "b", "a", want{nil, ""}},
		{"test 2", "c", "b", want{nil, ""}},
		{"test 3", "b", "a", want{e.ErrRepoExists, ""}},
	}

	for _, test := range addURLTests {
		t.Run(test.name, func(t *testing.T) {
			err := repo.AddURL(test.shortURL, test.longURL)
			if !errors.Is(err, test.want.err) {
				t.Errorf("AddURL test failed: shortURL %s, longURL %s, err: %v", test.shortURL, test.longURL, err)
			}
		})
	}
}

func testGetURL(t *testing.T, repo repository.URLRepository) {
	t.Helper()

	getURLTests := []test{
		{"test 1", "b", "a", want{nil, "b"}},
		{"test 2", "c", "b", want{nil, "c"}},
		{"test 3", "d", "c", want{e.ErrRepoNotFound, ""}},
	}

	for _, test := range getURLTests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			longURL, err := repo.GetURL(test.shortURL)
			if !errors.Is(err, test.want.err) || longURL != test.want.longURL {
				t.Errorf(
					"GetURL test failed: shortURL %s, expected longURL %s, got longURL %s, err: %v",
					test.shortURL, test.want.longURL, longURL, err,
				)
			}
		})
	}
}

func testMemento(t *testing.T, repo repository.URLRepository) {
	t.Helper()

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

func RunRepositoryTests(t *testing.T, repo repository.URLRepository) {
	t.Helper()
	testAddURL(t, repo)
	testGetURL(t, repo)
	testMemento(t, repo)
}

func TestInMemoryURLRepository(t *testing.T) {
	t.Parallel()

	repo := repository.NewInMemoryURLRepository()

	RunRepositoryTests(t, repo)
}
