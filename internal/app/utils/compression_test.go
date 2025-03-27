package utils_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	e "github.com/patraden/ya-practicum-go-shortly/internal/app/domain/errors"
	"github.com/patraden/ya-practicum-go-shortly/internal/app/utils"
)

func TestCompressDecompress(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		encoding    string
		data        []byte
		expectedErr error
	}{
		{
			name:        "Compress and Decompress Gzip",
			encoding:    "gzip",
			data:        []byte("Test data for gzip compression"),
			expectedErr: nil,
		},
		{
			name:        "Compress and Decompress Deflate",
			encoding:    "deflate",
			data:        []byte("Test data for deflate compression"),
			expectedErr: nil,
		},
		{
			name:        "Unsupported encoding",
			encoding:    "unsupported",
			data:        []byte("Test data for unsupported compression"),
			expectedErr: e.ErrUtilsCompEncoding,
		},
		{
			name:        "Empty Data",
			encoding:    "gzip",
			data:        []byte(""),
			expectedErr: nil,
		},
		{
			name:        "No compression (empty encoding string)",
			encoding:    "",
			data:        []byte("Test data with no compression"),
			expectedErr: nil,
		},
	}

	for _, tcase := range tests {
		t.Run(tcase.name, func(t *testing.T) {
			t.Parallel()

			compressedData, err := utils.Compress(tcase.data, tcase.encoding)

			if tcase.expectedErr != nil {
				assert.Equal(t, tcase.expectedErr, err)
				return // Skip decompression if compression failed
			}

			if tcase.encoding != "" {
				assert.NotEmpty(t, compressedData)
			}

			decompressedData, err := utils.Decompress(compressedData, tcase.encoding)

			require.NoError(t, err)
			assert.Equal(t, tcase.data, decompressedData)
		})
	}
}
