package utils

import (
	"bytes"
	"compress/flate"
	"compress/gzip"
	"io"

	e "github.com/patraden/ya-practicum-go-shortly/internal/app/domain/errors"
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
		return data, e.ErrUtilsCompEncoding
	}

	w := encoder(&buf, flate.BestCompression)
	if w == nil {
		return data, e.ErrUtilsEncoderOpen
	}

	wc, ok := w.(io.WriteCloser)
	if !ok {
		return data, e.ErrUtilsEncoderCast
	}

	_, err := wc.Write(data)
	if err != nil {
		return nil, e.Wrap("compression encoder write error", err, errLabel)
	}

	err = wc.Close()
	if err != nil {
		return nil, e.Wrap("compression encoder close error", err, errLabel)
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
		return data, e.ErrUtilsDecompionEncoding
	}

	r := decoder(bytes.NewReader(data))
	defer r.Close()

	var b bytes.Buffer

	_, err := b.ReadFrom(r)
	if err != nil {
		return nil, e.Wrap("decompression decoder read error", err, errLabel)
	}

	return b.Bytes(), nil
}
