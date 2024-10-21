package middleware

import (
	"compress/gzip"
	"compress/zlib"
	"io"
	"net/http"

	e "github.com/patraden/ya-practicum-go-shortly/internal/app/errors"
	"github.com/patraden/ya-practicum-go-shortly/internal/app/logger"
)

func Decompress(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var reader io.ReadCloser
		var err error

		encoding := r.Header.Get("Content-Encoding")
		switch encoding {
		case "gzip":
			reader, err = gzip.NewReader(r.Body)
		case "deflate":
			reader, err = zlib.NewReader(r.Body)
		case "":
			reader, err = r.Body, nil
		}

		if err != nil {
			logger.Log.Error().Msg(err.Error())
			http.Error(w, e.ErrDecompress.Error(), http.StatusBadRequest)
			return
		}

		if reader == nil {
			logger.Log.Error().Msg("decompression not implremeneted")
			http.Error(w, "decompression not implremeneted", http.StatusInternalServerError)
			return
		}

		defer reader.Close()
		r.Body = reader

		next.ServeHTTP(w, r)

	})
}