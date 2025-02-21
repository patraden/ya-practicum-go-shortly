package utils

import (
	"bytes"
	"compress/flate"
	"compress/gzip"
	"io"

	e "github.com/patraden/ya-practicum-go-shortly/internal/app/domain/errors"
)

// EncoderGzip returns a Gzip encoder with the specified compression level
// for writing compressed data to the provided writer.
func EncoderGzip(w io.Writer, level int) io.Writer {
	gw, err := gzip.NewWriterLevel(w, level)
	if err != nil {
		return nil
	}

	return gw
}

// EncoderDeflate returns a Deflate encoder with the specified compression level
// for writing compressed data to the provided writer.
func EncoderDeflate(w io.Writer, level int) io.Writer {
	dw, err := flate.NewWriter(w, level)
	if err != nil {
		return nil
	}

	return dw
}

// DecoderGzip returns a Gzip decompressor for reading compressed data from the provided reader.
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

// DecoderDeflate returns a Deflate decompressor for reading compressed data from the provided reader.
func DecoderDeflate(r io.Reader) io.ReadCloser {
	return flate.NewReader(r)
}

// Compress compresses the given data according to the specified encoding format (either "deflate" or "gzip").
// It returns the compressed data or an error if the compression process fails.
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

// Decompress decompresses the given data according to the specified encoding format (either "deflate" or "gzip").
// It returns the decompressed data or an error if the decompression process fails.
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
