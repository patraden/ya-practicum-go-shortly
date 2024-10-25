package repository_test

import (
	"testing"

	e "github.com/patraden/ya-practicum-go-shortly/internal/app/errors"
	"github.com/patraden/ya-practicum-go-shortly/internal/app/repository"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewInMemoryURLRepository(t *testing.T) {
	t.Parallel()

	storage := repository.NewInMemoryURLRepository()
	err := storage.AddURL("e", "f")
	require.NoError(t, err, "expected no error when adding a URL")

	type want struct {
		err error
	}

	tests := []struct {
		name  string
		key   string
		value string
		want  want
	}{
		{
			name:  "test 1",
			key:   "a",
			value: "b",
			want: want{
				err: nil,
			},
		},
		{
			name:  "test 2",
			key:   "b",
			value: "d",
			want: want{
				err: nil,
			},
		},
		{
			name:  "test 3",
			key:   "e",
			value: "f",
			want: want{
				err: e.ErrExists,
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			err := storage.AddURL(test.key, test.value)

			if test.want.err != nil {
				require.Error(t, err)
				assert.EqualError(t, err, test.want.err.Error())
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
