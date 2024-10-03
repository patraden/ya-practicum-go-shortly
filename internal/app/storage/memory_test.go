package storage

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMapKVStorageAdd(t *testing.T) {
	storage := NewMapStorage()

	type want struct {
		val string
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
				val: "a",
				err: nil,
			},
		},
		{
			name:  "test 2",
			key:   "b",
			value: "d",
			want: want{
				val: "b",
				err: nil,
			},
		},
		{
			name:  "test 3",
			key:   "b",
			value: "e",
			want: want{
				val: "b",
				err: ErrKeyExists,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v, err := storage.Add(tt.key, tt.value)
			if tt.want.err != nil {
				assert.Error(t, err)
				assert.EqualError(t, err, tt.want.err.Error())
			} else {
				assert.NoError(t, err)
				assert.Equal(t, v, tt.want.val)
				assert.Equal(t, v, tt.want.val)
			}
		})
	}
}
