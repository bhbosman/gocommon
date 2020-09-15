package stream

import "io"

type ByteReaderCloser interface {
	io.Closer
	io.ByteReader
}

type byteReaderNoCloser struct {
	byteReader io.ByteReader
}

func NewByteReaderNoCloser(byteReader io.ByteReader) *byteReaderNoCloser {
	return &byteReaderNoCloser{
		byteReader: byteReader,
	}
}

func (b byteReaderNoCloser) Close() error {
	return nil
}

func (b *byteReaderNoCloser) ReadByte() (byte, error) {
	return b.byteReader.ReadByte()
}

type ByteReaderWithCloser struct {
	byteReader io.ByteReader
	closer     io.Closer
}

func (b ByteReaderWithCloser) ReadByte() (byte, error) {
	return b.byteReader.ReadByte()
}

func NewByteReaderWithCloser(closer io.Closer, byteReader io.ByteReader) *ByteReaderWithCloser {
	return &ByteReaderWithCloser{
		byteReader: byteReader,
		closer:     closer,
	}
}

func (b ByteReaderWithCloser) Close() error {
	return b.closer.Close()
}

type NullWriter struct {
}

func (self NullWriter) Write(p []byte) (n int, err error) {
	return len(p), nil
}
