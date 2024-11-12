package repository_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/patraden/ya-practicum-go-shortly/internal/app/domain"
	e "github.com/patraden/ya-practicum-go-shortly/internal/app/domain/errors"
	"github.com/patraden/ya-practicum-go-shortly/internal/app/repository"
)

type test struct {
	name    string
	mapping *domain.URLMapping
	wantErr error
}

func TestAddURL(t *testing.T) {
	t.Parallel()

	repo := repository.NewInMemoryURLRepository()

	addURLTests := []test{
		{"unique slug and URL", domain.NewURLMapping("slug1", "url1"), nil},
		{"duplicate URL with different slug", domain.NewURLMapping("slug2", "url1"), e.ErrOriginalExists},
		{"duplicate slug", domain.NewURLMapping("slug1", "url2"), e.ErrSlugExists},
	}

	for _, tc := range addURLTests {
		t.Run(tc.name, func(t *testing.T) {
			_, err := repo.AddURLMapping(context.Background(), tc.mapping)
			require.ErrorIs(t, err, tc.wantErr)
		})
	}
}

func TestGetURL(t *testing.T) {
	t.Parallel()

	repo := repository.NewInMemoryURLRepository()

	_, err := repo.AddURLMapping(context.Background(), domain.NewURLMapping("slug1", "url1"))
	require.NoError(t, err)

	_, err = repo.AddURLMapping(context.Background(), domain.NewURLMapping("slug2", "url2"))
	require.NoError(t, err)

	getURLTests := []test{
		{"existing slug 'slug1'", domain.NewURLMapping("slug1", "url1"), nil},
		{"existing slug 'slug2'", domain.NewURLMapping("slug2", "url2"), nil},
		{"nonexistent slug 'slug3'", domain.NewURLMapping("slug3", "url3"), e.ErrSlugNotFound},
	}

	for _, tc := range getURLTests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			_, err := repo.GetURLMapping(context.Background(), tc.mapping.Slug)
			if tc.wantErr == nil {
				require.NoError(t, err)
			} else {
				require.ErrorIs(t, err, tc.wantErr)
			}
		})
	}
}

func TestAddURLBatch(t *testing.T) {
	t.Parallel()

	repo := repository.NewInMemoryURLRepository()

	batch := &[]domain.URLMapping{
		*domain.NewURLMapping("slug1", "url1"),
		*domain.NewURLMapping("slug2", "url2"),
		*domain.NewURLMapping("slug3", "url3"),
	}

	err := repo.AddURLMappingBatch(context.Background(), batch)
	require.NoError(t, err)

	duplicateSlugBatch := &[]domain.URLMapping{
		*domain.NewURLMapping("slug1", "url5"),
		*domain.NewURLMapping("slug4", "url4"),
	}
	err = repo.AddURLMappingBatch(context.Background(), duplicateSlugBatch)
	require.ErrorIs(t, err, e.ErrSlugExists)

	duplicateURLBatch := &[]domain.URLMapping{
		*domain.NewURLMapping("slug5", "url1"),
		*domain.NewURLMapping("slug6", "url6"),
	}
	err = repo.AddURLMappingBatch(context.Background(), duplicateURLBatch)
	require.ErrorIs(t, err, e.ErrOriginalExists)
}

func TestMementoOps(t *testing.T) {
	t.Parallel()

	repo := repository.NewInMemoryURLRepository()
	ctx := context.Background()

	_, err := repo.AddURLMapping(ctx, domain.NewURLMapping("slug1", "url1"))
	require.NoError(t, err)

	_, err = repo.AddURLMapping(ctx, domain.NewURLMapping("slug2", "url2"))
	require.NoError(t, err)

	initialMemento, err := repo.CreateMemento()
	require.NoError(t, err)

	_, err = repo.AddURLMapping(ctx, domain.NewURLMapping("slug3", "url3"))
	require.NoError(t, err)

	_, err = repo.GetURLMapping(ctx, "slug1")
	require.NoError(t, err)

	err = repo.RestoreMemento(initialMemento)
	require.NoError(t, err)

	res, err := repo.GetURLMapping(ctx, "slug1")
	require.NoError(t, err)
	require.Equal(t, "url1", string(res.OriginalURL))

	_, err = repo.GetURLMapping(ctx, "slug3")
	require.ErrorIs(t, err, e.ErrSlugNotFound)

	err = repo.RestoreMemento(nil)
	require.NoError(t, err)
}
