package file_test

import (
	"io"
	"os"
	"testing"

	e "github.com/patraden/ya-practicum-go-shortly/internal/app/errors"
	"github.com/patraden/ya-practicum-go-shortly/internal/app/logger"
	"github.com/patraden/ya-practicum-go-shortly/internal/app/repository/file"
	"github.com/patraden/ya-practicum-go-shortly/internal/app/utils"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFileReadWrite(t *testing.T) {
	t.Parallel()

	testRecords := []*file.Record{
		{
			UUID:     utils.UUID(),
			ShortURL: "XXXYYYZZ",
			LongURL:  "http://ya.ru",
		},
		{
			UUID:     utils.UUID(),
			ShortURL: "XXXYYZZZ",
			LongURL:  "http://ya.com",
		},
		{
			UUID:     utils.UUID(),
			ShortURL: "XXXYYZZZ",
			LongURL:  "http://ya.com",
		},
	}

	t.Run("test writer and reader", func(t *testing.T) {
		t.Parallel()

		log := logger.NewLogger(zerolog.InfoLevel).GetLogger()

		fileName := "records.json"
		defer os.Remove(fileName)

		writer := file.NewWriter(fileName, log)
		defer writer.Close()

		reader := file.NewReader(fileName, log)
		defer reader.Close()

		for _, record := range testRecords {
			err := writer.Write(record)
			require.NoError(t, err)
		}

		cache := make(map[string]string)
		err := reader.LoadToCache(&cache)
		require.NoError(t, err)
		assert.Equal(t, len(testRecords), len(cache)+1)

		_, err = reader.Find("XXXYYZZZ")
		require.NoError(t, err)

		_, err = reader.Find("YYYY")
		require.ErrorIs(t, err, e.ErrRepoFile)

		err = reader.Reset()
		require.NoError(t, err)

		_, err = reader.Read()
		require.NoError(t, err)

		_, err = reader.Read()
		require.NoError(t, err)

		_, err = reader.Read()
		require.NoError(t, err)

		_, err = reader.Read()
		require.ErrorIs(t, err, io.EOF)
	})
}
