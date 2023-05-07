package gocommon

import (
	"github.com/reactivex/rxgo/v2"
	"io"
)

type IObservable interface {
	Observe(opts ...rxgo.Option) <-chan rxgo.Item
	Map(apply rxgo.Func, opts ...rxgo.Option) rxgo.Observable
	FlatMap(apply rxgo.ItemToObservable, opts ...rxgo.Option) rxgo.Observable
	//rxgo.Observable
}

type ByteReaderCloser interface {
	io.ByteReader
	io.Closer
}
type byteReaderCloser struct {
	br io.ByteReader
}

func (self *byteReaderCloser) ReadByte() (byte, error) {
	return self.br.ReadByte()
}

func (self *byteReaderCloser) Close() error {
	return nil
}

func NewByteReaderNoCloser(br io.ByteReader) ByteReaderCloser {
	return &byteReaderCloser{
		br: br,
	}
}

type byteReaderWithCloser struct {
	br io.ByteReader
	c  io.Closer
}

func (self *byteReaderWithCloser) Close() error {
	return self.c.Close()
}

func (self *byteReaderWithCloser) ReadByte() (byte, error) {
	return self.br.ReadByte()
}

func NewByteReaderWithCloser(c io.Closer, br io.ByteReader) ByteReaderCloser {
	return &byteReaderWithCloser{
		br: br,
		c:  c,
	}
}
