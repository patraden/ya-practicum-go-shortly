package utils

import (
	"bytes"
	"compress/flate"
	"compress/gzip"
	"fmt"
	"io"

	e "github.com/patraden/ya-practicum-go-shortly/internal/app/errors"
)

func EncoderGzip(w io.Writer, level int) io.Writer {
	gw, err := gzip.NewWriterLevel(w, level)
	if err != nil {
		return nil
	}

	return gw
}

func EncoderDeflate(w io.Writer, level int) io.Writer {
	dw, err := flate.NewWriter(w, level)
	if err != nil {
		return nil
	}

	return dw
}

func DecoderGzip(r io.Reader) io.ReadCloser {
	if r == nil {
		return io.NopCloser(bytes.NewReader([]byte{}))
	}

	dr, err := gzip.NewReader(r)
	if err != nil {
		return nil
	}

	return dr
}

func DecoderDeflate(r io.Reader) io.ReadCloser {
	return flate.NewReader(r)
}

// helper for compression tests.
func Compress(data []byte, encoding string) ([]byte, error) {
	var buf bytes.Buffer
	var encoder func(w io.Writer, level int) io.Writer

	switch encoding {
	case "deflate":
		encoder = EncoderDeflate
	case "gzip":
		encoder = EncoderGzip
	case "":
		return data, nil
	default:
		return data, fmt.Errorf("%w: bad encoding %s", e.ErrUtils, encoding)
	}

	w := encoder(&buf, flate.BestCompression)
	if w == nil {
		return data, fmt.Errorf("%w: could't get writer for %s", e.ErrUtils, encoding)
	}

	wc, ok := w.(io.WriteCloser)
	if !ok {
		return data, fmt.Errorf("%w: could't get writecloser for %s", e.ErrUtils, encoding)
	}

	_, err := wc.Write(data)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", e.ErrUtils.Error(), err)
	}

	err = wc.Close()
	if err != nil {
		return nil, fmt.Errorf("%s: %w", e.ErrUtils.Error(), err)
	}

	return buf.Bytes(), nil
}

// helper for compression tests.
func Decompress(data []byte, encoding string) ([]byte, error) {
	var decoder func(r io.Reader) io.ReadCloser

	switch encoding {
	case "deflate":
		decoder = DecoderDeflate
	case "gzip":
		decoder = DecoderGzip
	case "":
		return data, nil
	default:
		return data, fmt.Errorf("%w: bad encoding %s", e.ErrUtils, encoding)
	}

	r := decoder(bytes.NewReader(data))
	defer r.Close()

	var b bytes.Buffer

	_, err := b.ReadFrom(r)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", e.ErrUtils.Error(), err)
	}

	return b.Bytes(), nil
}
