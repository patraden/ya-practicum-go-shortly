// Decompressor approach is inspired
// by https://github.com/go-chi/chi/blob/master/middleware/compress.go.
// and implemneted in a simplistic form for now.

package middleware

import (
	"io"
	"net/http"
	"strings"
	"sync"

	"github.com/patraden/ya-practicum-go-shortly/internal/app/utils"
)

func Decompress() func(next http.Handler) http.Handler {
	decompressor := NewDecompressor()

	return decompressor.Handler
}

type DecodeFunc func(io.Reader) io.ReadCloser

type Decompressor struct {
	pooledDecoders map[string]*sync.Pool
}

func NewDecompressor() *Decompressor {
	d := &Decompressor{
		pooledDecoders: make(map[string]*sync.Pool),
	}

	d.SetDecoder("deflate", utils.DecoderDeflate)
	d.SetDecoder("gzip", utils.DecoderGzip)

	return d
}

func (d *Decompressor) SetDecoder(encoding string, fn DecodeFunc) {
	encoding = strings.ToLower(encoding)

	delete(d.pooledDecoders, encoding)

	if fn(nil) != nil {
		pool := &sync.Pool{
			New: func() interface{} {
				return fn
			},
		}
		d.pooledDecoders[encoding] = pool
	}
}

func (d *Decompressor) selectDecoder(h http.Header, r io.ReadCloser) (io.ReadCloser, string) {
	encoded := h.Get("Content-Encoding")

	// content is not encoded
	if encoded == "" {
		return r, ""
	}

	// try to get from pooledDecoders
	for name := range d.pooledDecoders {
		if name == encoded {
			if pool, ok := d.pooledDecoders[name]; ok {
				if decoder, ok := pool.Get().(DecodeFunc); ok {
					return decoder(r), encoded
				}
			}
		}
	}

	return nil, encoded
}

func (d *Decompressor) Handler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		decoder, _ := d.selectDecoder(r.Header, r.Body)

		if decoder == nil {
			http.Error(w, "Content could not be decoded", http.StatusInternalServerError)

			return
		}

		r.Body = decoder
		defer decoder.Close()
		next.ServeHTTP(w, r)
	})
}
