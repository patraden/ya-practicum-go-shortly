package repository

import (
	"testing"

	e "github.com/patraden/ya-practicum-go-shortly/internal/app/errors"
	"github.com/stretchr/testify/assert"
)

func TestNewInMemoryURLRepository(t *testing.T) {
	storage := NewInMemoryURLRepository()

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
			key:   "b",
			value: "e",
			want: want{
				err: e.ErrExists,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := storage.AddURL(tt.key, tt.value)
			if tt.want.err != nil {
				assert.Error(t, err)
				assert.EqualError(t, err, tt.want.err.Error())
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
