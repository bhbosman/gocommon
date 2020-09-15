package multiBlock

import (
	"bytes"
	"crypto/sha1"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestReaderWriter(t *testing.T) {
	t.Run("empty read", func(t *testing.T) {
		rw := NewReaderWriter()
		read, err := rw.Read([]byte{0})
		assert.Error(t, err)
		assert.Equal(t, 0, read)
	})

	t.Run("one 4096W, 4 x 1024R, empty buffer", func(t *testing.T) {
		rw := NewReaderWriter()
		n, err := rw.Write(make([]byte, 4096))
		assert.NoError(t, err)
		assert.Equal(t, 4096, n)
		readBuffer := make([]byte, 1024)
		for i := 0; i <= 3; i++ {
			n, err = rw.Read(readBuffer)
			assert.NoError(t, err)
			assert.Equal(t, 1024, n)
		}
		assert.Equal(t, 0, rw.Size())
		n, err = rw.Read(readBuffer)
		assert.Error(t, err)
		assert.Equal(t, 0, n)
	})

	t.Run("one 4096W, 4 x 1024R, empty buffer", func(t *testing.T) {
		rw := NewReaderWriter()
		for i := 0; i < 8; i++ {
			n, err := rw.Write(make([]byte, 1024))
			assert.NoError(t, err)
			assert.Equal(t, 1024, n)
		}
		readBuffer := make([]byte, 2048)
		for i := 0; i < 4; i++ {
			n, err := rw.Read(readBuffer)
			assert.NoError(t, err)
			assert.Equal(t, 2048, n)
		}
		assert.Equal(t, 0, rw.Size())
		n, err := rw.Read(readBuffer)
		assert.Error(t, err)
		assert.Equal(t, 0, n)
	})

	t.Run("", func(t *testing.T) {
		rw := NewReaderWriterSize(512)
		for i := 0; i < 8; i++ {
			n, err := rw.Write(make([]byte, 1024))
			assert.NoError(t, err)
			assert.Equal(t, 1024, n)
		}
		rw.Flatten()
	})

	t.Run("flatten to next", func(t *testing.T) {
		rw01 := NewReaderWriterSize(512)
		_, _ = rw01.Write([]byte{1})
		rw02 := NewReaderWriterSize(512)
		_, _ = rw02.Write([]byte{2})
		rw01.Add(rw02)
		assert.Equal(t, 2, rw01.Size())
		assert.Equal(t, 1022, rw01.Waste())
		flatten, _ := rw01.Flatten()
		assert.Len(t, flatten, 2)
		assert.Equal(t, 2, rw01.Size())
		assert.Equal(t, 510, rw01.Waste())
	})

	t.Run("flatten to next", func(t *testing.T) {
		rw01 := NewReaderWriterSize(128)
		_, _ = rw01.Write(make([]byte, 20))
		_, _ = rw01.Write(make([]byte, 2048))
		assert.Equal(t, 2068, rw01.Size())
		assert.Equal(t, 108, rw01.Waste())

		rw02 := NewReaderWriterSize(256)
		_, _ = rw02.Write(make([]byte, 20))
		_, _ = rw02.Write(make([]byte, 4096))
		assert.Equal(t, 4116, rw02.Size())
		assert.Equal(t, 236, rw02.Waste())

		rw01.Add(rw02)
		assert.Equal(t, 4116+2068, rw01.Size())
		assert.Equal(t, 108+236, rw01.Waste())
		flatten, _ := rw01.Flatten()
		assert.Len(t, flatten, 4116+2068)
		assert.Equal(t, 4116+2068, rw01.Size())
		assert.Equal(t, 0, rw01.Waste())
	})
	t.Run("Add Three Buffers, verify with hash", func(t *testing.T) {
		data001 := []byte("ReaderWriter_test_0001")
		data002 := []byte("ReaderWriter_test_0002")
		data003 := []byte("ReaderWriter_test_0003")
		sha1Hash := sha1.New()
		_, err := sha1Hash.Write(data001)
		assert.NoError(t, err)
		_, err = sha1Hash.Write(data002)
		assert.NoError(t, err)
		_, err = sha1Hash.Write(data003)
		assert.NoError(t, err)

		mustbe := sha1Hash.Sum(nil)
		sha1Hash.Reset()

		rw01 := NewReaderWriterSize(128)
		_, err = rw01.Write(data001)
		assert.NoError(t, err)
		rw02 := NewReaderWriterSize(128)
		_, err = rw02.Write(data002)
		assert.NoError(t, err)

		rw03 := NewReaderWriterSize(128)
		_, err = rw03.Write(data003)
		assert.NoError(t, err)

		err = rw01.Add(rw02, rw03)
		assert.NoError(t, err)

		err = rw01.Dump(sha1Hash)
		assert.NoError(t, err)

		valueIs := sha1Hash.Sum(nil)
		t.Logf("%v\n", mustbe)
		t.Logf("%v\n", valueIs)
		assert.Equal(t, mustbe, valueIs)
	})
	t.Run("Hash", func(t *testing.T) {
		t.Run("EmptyHash", func(t *testing.T) {
			sha1Hash := sha1.New()
			mustbe := sha1Hash.Sum(nil)
			sha1Hash.Reset()
			rw01 := NewReaderWriterSize(128)
			err := rw01.Dump(sha1Hash)
			assert.NoError(t, err)

			valueIs := sha1Hash.Sum(nil)
			t.Logf("%v\n", mustbe)
			t.Logf("%v\n", valueIs)
			assert.Equal(t, mustbe, valueIs)
		})

		t.Run("One Buffer", func(t *testing.T) {
			data := []byte("ReaderWriter_test")
			sha1Hash := sha1.New()
			_, err := sha1Hash.Write(data)
			assert.NoError(t, err)
			mustbe := sha1Hash.Sum(nil)
			sha1Hash.Reset()

			rw01 := NewReaderWriterSize(128)
			_, err = rw01.Write(data)
			assert.NoError(t, err)
			err = rw01.Dump(sha1Hash)
			assert.NoError(t, err)

			valueIs := sha1Hash.Sum(nil)
			t.Logf("%v\n", mustbe)
			t.Logf("%v\n", valueIs)
			assert.Equal(t, mustbe, valueIs)
		})

		t.Run("Two Buffers", func(t *testing.T) {
			data001 := []byte("ReaderWriter_test_0001")
			data002 := []byte("ReaderWriter_test_0002")
			sha1Hash := sha1.New()
			_, err := sha1Hash.Write(data001)
			assert.NoError(t, err)
			_, err = sha1Hash.Write(data002)
			assert.NoError(t, err)

			mustbe := sha1Hash.Sum(nil)
			sha1Hash.Reset()

			rw01 := NewReaderWriterSize(128)
			_, err = rw01.Write(data001)
			assert.NoError(t, err)
			rw02 := NewReaderWriterSize(128)
			_, err = rw01.Write(data002)
			assert.NoError(t, err)

			err = rw01.Add(rw02)
			assert.NoError(t, err)

			err = rw01.Dump(sha1Hash)
			assert.NoError(t, err)

			valueIs := sha1Hash.Sum(nil)
			t.Logf("%v\n", mustbe)
			t.Logf("%v\n", valueIs)
			assert.Equal(t, mustbe, valueIs)
		})

		t.Run("Three Buffers", func(t *testing.T) {
			data001 := []byte("ReaderWriter_test_0001")
			data002 := []byte("ReaderWriter_test_0002")
			data003 := []byte("ReaderWriter_test_0003")
			sha1Hash := sha1.New()
			_, err := sha1Hash.Write(data001)
			assert.NoError(t, err)
			_, err = sha1Hash.Write(data002)
			assert.NoError(t, err)
			_, err = sha1Hash.Write(data003)
			assert.NoError(t, err)

			mustbe := sha1Hash.Sum(nil)
			sha1Hash.Reset()

			rw01 := NewReaderWriterSize(128)
			_, err = rw01.Write(data001)
			assert.NoError(t, err)
			rw02 := NewReaderWriterSize(128)
			_, err = rw02.Write(data002)
			assert.NoError(t, err)

			rw03 := NewReaderWriterSize(128)
			_, err = rw03.Write(data003)
			assert.NoError(t, err)

			err = rw01.Add(rw02)
			assert.NoError(t, err)

			err = rw01.Add(rw03)
			assert.NoError(t, err)

			err = rw01.Dump(sha1Hash)
			assert.NoError(t, err)

			valueIs := sha1Hash.Sum(nil)
			t.Logf("%v\n", mustbe)
			t.Logf("%v\n", valueIs)
			assert.Equal(t, mustbe, valueIs)
		})

	})
	t.Run("Three io.ReadeWriter Buffers", func(t *testing.T) {
		data001 := []byte("ReaderWriter_test_0001")
		data002 := []byte("ReaderWriter_test_0002")
		data003 := []byte("ReaderWriter_test_0003")
		sha1Hash := sha1.New()
		_, err := sha1Hash.Write(data001)
		assert.NoError(t, err)
		_, err = sha1Hash.Write(data002)
		assert.NoError(t, err)
		_, err = sha1Hash.Write(data003)
		assert.NoError(t, err)

		mustbe := sha1Hash.Sum(nil)
		sha1Hash.Reset()

		rw01 := bytes.Buffer{}
		_, err = rw01.Write(data001)
		assert.NoError(t, err)

		rw02 := bytes.Buffer{}
		_, err = rw02.Write(data002)
		assert.NoError(t, err)

		rw03 := bytes.Buffer{}
		_, err = rw03.Write(data003)
		assert.NoError(t, err)

		ma := NewReaderWriter()
		err = ma.AddReaders(&rw01, &rw02, &rw03)
		if !assert.NoError(t, err) {
			return
		}

		err = ma.Dump(sha1Hash)
		if !assert.NoError(t, err) {
			return
		}

		valueIs := sha1Hash.Sum(nil)
		t.Logf("%v\n", mustbe)
		t.Logf("%v\n", valueIs)
		assert.Equal(t, mustbe, valueIs)
	})
	t.Run("setnext", func(t *testing.T) {
		b1 := NewReaderWriterBlock([]byte{1, 2, 3, 4, 5})
		b2 := NewReaderWriterBlock([]byte{6, 7, 8, 9, 10})
		b1.SetNext(b2)
		data := make([]byte, 10)
		n, _ := b1.Read(data)
		assert.Equal(t, 10, n)
		assert.Equal(t, []byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}, data)
	})
	t.Run("setnext", func(t *testing.T) {
		b1 := NewReaderWriterBlock([]byte{1, 2, 3, 4, 5})
		b2 := NewReaderWriterBlock([]byte{6, 7, 8, 9, 10})
		b3 := NewReaderWriterBlock([]byte{11, 12, 13, 14, 15})
		_ = b1.SetNext(b2)
		_ = b1.SetNext(b3)
		data := make([]byte, 15)
		n, _ := b1.Read(data)
		assert.Equal(t, 15, n)
		assert.Equal(t, []byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15}, data)
	})
}
