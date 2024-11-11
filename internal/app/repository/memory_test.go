package repository_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/patraden/ya-practicum-go-shortly/internal/app/domain"
	e "github.com/patraden/ya-practicum-go-shortly/internal/app/errors"
	"github.com/patraden/ya-practicum-go-shortly/internal/app/repository"
)

type test struct {
	name    string
	mapping *domain.URLMapping
	wantErr error
}

func testAddURL(t *testing.T, repo repository.URLRepository) {
	t.Helper()

	addURLTests := []test{
		{"test 1", domain.NewURLMapping("b", "a"), nil},
		{"test 2", domain.NewURLMapping("c", "b"), nil},
		{"test 3", domain.NewURLMapping("b", "a"), e.ErrRepoExists},
	}

	for _, test := range addURLTests {
		t.Run(test.name, func(t *testing.T) {
			err := repo.AddURLMapping(context.Background(), test.mapping)
			require.ErrorIs(t, err, test.wantErr)
		})
	}
}

func testGetURL(t *testing.T, repo repository.URLRepository) {
	t.Helper()

	getURLTests := []test{
		{"test 1", domain.NewURLMapping("b", "a"), nil},
		{"test 2", domain.NewURLMapping("c", "b"), nil},
		{"test 3", domain.NewURLMapping("d", "c"), e.ErrRepoNotFound},
	}

	for _, test := range getURLTests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			_, err := repo.GetURLMapping(context.Background(), test.mapping.Slug)
			require.ErrorIs(t, err, test.wantErr)
		})
	}
}

func TestInMemoryURLRepository(t *testing.T) {
	t.Parallel()

	repo := repository.NewInMemoryURLRepository()
	testAddURL(t, repo)
	testGetURL(t, repo)
}
