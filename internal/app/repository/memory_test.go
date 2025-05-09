package repository_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/patraden/ya-practicum-go-shortly/internal/app/domain"
	e "github.com/patraden/ya-practicum-go-shortly/internal/app/domain/errors"
	"github.com/patraden/ya-practicum-go-shortly/internal/app/dto"
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
	userID := domain.NewUserID()

	addURLTests := []test{
		{"unique slug and URL", domain.NewURLMapping("slug1", "url1", userID), nil},
		{"duplicate URL with different slug", domain.NewURLMapping("slug2", "url1", userID), e.ErrOriginalExists},
		{"duplicate slug", domain.NewURLMapping("slug1", "url2", userID), e.ErrSlugExists},
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
	userID := domain.NewUserID()

	_, err := repo.AddURLMapping(context.Background(), domain.NewURLMapping("slug1", "url1", userID))
	require.NoError(t, err)

	_, err = repo.AddURLMapping(context.Background(), domain.NewURLMapping("slug2", "url2", userID))
	require.NoError(t, err)

	getURLTests := []test{
		{"existing slug 'slug1'", domain.NewURLMapping("slug1", "url1", userID), nil},
		{"existing slug 'slug2'", domain.NewURLMapping("slug2", "url2", userID), nil},
		{"nonexistent slug 'slug3'", domain.NewURLMapping("slug3", "url3", userID), e.ErrSlugNotFound},
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
	userID := domain.NewUserID()

	batch := &[]domain.URLMapping{
		*domain.NewURLMapping("slug1", "url1", userID),
		*domain.NewURLMapping("slug2", "url2", userID),
		*domain.NewURLMapping("slug3", "url3", userID),
	}

	err := repo.AddURLMappingBatch(context.Background(), batch)
	require.NoError(t, err)

	duplicateSlugBatch := &[]domain.URLMapping{
		*domain.NewURLMapping("slug1", "url5", userID),
		*domain.NewURLMapping("slug4", "url4", userID),
	}
	err = repo.AddURLMappingBatch(context.Background(), duplicateSlugBatch)
	require.ErrorIs(t, err, e.ErrSlugExists)

	duplicateURLBatch := &[]domain.URLMapping{
		*domain.NewURLMapping("slug5", "url1", userID),
		*domain.NewURLMapping("slug6", "url6", userID),
	}
	err = repo.AddURLMappingBatch(context.Background(), duplicateURLBatch)
	require.ErrorIs(t, err, e.ErrOriginalExists)
}

func TestMementoOps(t *testing.T) {
	t.Parallel()

	repo := repository.NewInMemoryURLRepository()
	ctx := context.Background()
	userID := domain.NewUserID()

	_, err := repo.AddURLMapping(ctx, domain.NewURLMapping("slug1", "url1", userID))
	require.NoError(t, err)

	_, err = repo.AddURLMapping(ctx, domain.NewURLMapping("slug2", "url2", userID))
	require.NoError(t, err)

	initialMemento, err := repo.CreateMemento()
	require.NoError(t, err)

	_, err = repo.AddURLMapping(ctx, domain.NewURLMapping("slug3", "url3", userID))
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

func TestMemGetUserURLMappings(t *testing.T) {
	t.Parallel()

	repo := repository.NewInMemoryURLRepository()
	userID := domain.NewUserID()
	otherUserID := domain.NewUserID()

	// Prepare data
	_, err := repo.AddURLMapping(context.Background(), domain.NewURLMapping("slug1", "url1", userID))
	require.NoError(t, err)

	_, err = repo.AddURLMapping(context.Background(), domain.NewURLMapping("slug2", "url2", userID))
	require.NoError(t, err)

	_, err = repo.AddURLMapping(context.Background(), domain.NewURLMapping("slug3", "url3", otherUserID))
	require.NoError(t, err)

	getUserURLMappingsTests := []struct {
		name      string
		userID    domain.UserID
		wantCount int
		wantErr   error
	}{
		{"existing user with mappings", userID, 2, nil},
		{"other user with mappings", otherUserID, 1, nil},
		{"nonexistent user", domain.NewUserID(), 0, e.ErrUserNotFound},
	}

	for _, tc := range getUserURLMappingsTests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			result, err := repo.GetUserURLMappings(context.Background(), tc.userID)

			if tc.wantErr == nil {
				require.NoError(t, err)
				require.Len(t, result, tc.wantCount)
			} else {
				require.ErrorIs(t, err, tc.wantErr)
				require.Nil(t, result)
			}
		})
	}
}

func TestDelUserURLMappings(t *testing.T) {
	t.Parallel()

	repo := repository.NewInMemoryURLRepository()
	userID := domain.NewUserID()
	otherUserID := domain.NewUserID()
	ctx := context.Background()

	_, err := repo.AddURLMapping(ctx, domain.NewURLMapping("slug1", "url1", userID))
	require.NoError(t, err)

	_, err = repo.AddURLMapping(ctx, domain.NewURLMapping("slug2", "url2", userID))
	require.NoError(t, err)

	_, err = repo.AddURLMapping(ctx, domain.NewURLMapping("slug3", "url3", otherUserID))
	require.NoError(t, err)

	tasks := []dto.UserSlug{
		{Slug: "slug1", UserID: userID},
		{Slug: "slug2", UserID: userID},
		{Slug: "slug4", UserID: otherUserID},
	}

	err = repo.DelUserURLMappings(ctx, tasks)
	require.NoError(t, err)

	m1, err := repo.GetURLMapping(ctx, "slug1")
	require.NoError(t, err)
	assert.True(t, m1.Deleted)

	m2, err := repo.GetURLMapping(ctx, "slug2")
	require.NoError(t, err)
	assert.True(t, m2.Deleted)

	m3, err := repo.GetURLMapping(ctx, "slug3")
	require.NoError(t, err)
	assert.False(t, m3.Deleted)
}

func TestGetStats(t *testing.T) {
	t.Parallel()

	repo := repository.NewInMemoryURLRepository()
	userID := domain.NewUserID()
	otherUserID := domain.NewUserID()
	ctx := context.Background()

	_, err := repo.AddURLMapping(ctx, domain.NewURLMapping("slug1", "url1", userID))
	require.NoError(t, err)

	_, err = repo.AddURLMapping(ctx, domain.NewURLMapping("slug2", "url2", userID))
	require.NoError(t, err)

	_, err = repo.AddURLMapping(ctx, domain.NewURLMapping("slug3", "url3", otherUserID))
	require.NoError(t, err)

	stats, err := repo.GetStats(ctx)
	require.NoError(t, err)

	assert.Equal(t, int64(2), stats.CountUsers)
	assert.Equal(t, int64(3), stats.CountSlugs)
}
