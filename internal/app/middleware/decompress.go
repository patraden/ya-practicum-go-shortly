// Package middleware provides the Decompress middleware to handle decompression
// of request bodies based on the "Content-Encoding" header.
// It supports decompressing "deflate" and "gzip" encodings using a pooled decoder approach.
// Approach is inspired by https://github.com/go-chi/chi/blob/master/middleware/compress.go.
package middleware

import (
	"io"
	"net/http"
	"strings"
	"sync"

	"github.com/patraden/ya-practicum-go-shortly/internal/app/utils"
)

// Decompress is a middleware that decompresses request bodies based on the
// "Content-Encoding" header. It uses a pooled decoder for supported encodings.
func Decompress() func(next http.Handler) http.Handler {
	decompressor := NewDecompressor()

	return decompressor.Handler
}

// DecodeFunc defines the function signature for decompressing an io.Reader.
type DecodeFunc func(io.Reader) io.ReadCloser

// Decompressor is a middleware that manages decompression for various encodings.
type Decompressor struct {
	pooledDecoders map[string]*sync.Pool
}

// NewDecompressor creates a new Decompressor instance, initializes the pool of decoders,
// and registers support for "deflate" and "gzip" encodings.
func NewDecompressor() *Decompressor {
	d := &Decompressor{
		pooledDecoders: make(map[string]*sync.Pool),
	}

	d.SetDecoder("deflate", utils.DecoderDeflate)
	d.SetDecoder("gzip", utils.DecoderGzip)

	return d
}

// SetDecoder registers a decoder function for a specific encoding type.
// If a decoder is already registered for the encoding, it is replaced.
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

// Handler is the HTTP middleware handler that decompresses the request body
// if the "Content-Encoding" header indicates compression. It updates the request's body
// to the decompressed version and passes it along to the next handler.
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
