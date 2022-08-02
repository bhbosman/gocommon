package stream

import (
	"encoding/binary"
	"fmt"
	"github.com/bhbosman/goerrors"
	"github.com/bhbosman/gomessageblock"
	"github.com/bhbosman/goprotoextra"
	"google.golang.org/protobuf/proto"
	"io"
	"reflect"
)

func Marshall(m proto.Message) (*gomessageblock.ReaderWriter, error) {
	if tc, ok := m.(interface {
		TypeCode() uint32
	}); ok {
		marshallBytes, err := proto.Marshal(m)
		if err != nil {
			return nil, err
		}
		tcBytes := [8]byte{}
		binary.LittleEndian.PutUint32(tcBytes[0:4], tc.TypeCode())
		binary.LittleEndian.PutUint32(tcBytes[4:8], uint32(len(marshallBytes)))
		return gomessageblock.NewReaderWriterWithBlocks(tcBytes[:], marshallBytes), nil
	}
	return nil, goerrors.NewInvalidNilParamError("No TypeCode() function")
}

func UnMarshal(rws goprotoextra.ReadWriterSize) (msgWrapper goprotoextra.IMessageWrapper, err error) {
	tc, err := rws.ReadTypeCode()
	if err != nil {
		return nil, err
	}
	if r, ok := registrationMap[tc]; ok {
		msg := r.CreateMessage()
		err = UnMarshalMessage(rws, msg)
		if err != nil {
			return nil, err
		}
		return r.CreateWrapper(msg)
	}
	return nil, goerrors.NewInvalidParamError(
		"tc",
		fmt.Sprintf("Invalid typecode %v", tc),
	)
}

func UnMarshalMessage(rws goprotoextra.ReadWriterSize, m proto.Message) error {
	if tc, ok := m.(interface {
		TypeCode() uint32
	}); ok {
		tcBytes := [8]byte{}
		_, err := rws.Read(tcBytes[:])
		if err != nil {
			return err
		}
		tcFromStream := binary.LittleEndian.Uint32(tcBytes[0:4])
		if tc.TypeCode() != tcFromStream {
			return goerrors.NewInvalidParamError(
				"tc.TypeCode",
				fmt.Sprintf(
					"Message %v type code does not match stream type code %v",
					reflect.TypeOf(m).String(),
					tcFromStream,
				),
			)
		}
		n := binary.LittleEndian.Uint32(tcBytes[4:8])
		if n == 0 {
			return nil
		}

		if n < uint32(rws.Size()) {
			return io.ErrShortBuffer
		}

		data := make([]byte, n)
		nn, err := rws.Read(data)
		if err != nil {
			return err
		}
		if uint32(nn) != n {
			return io.ErrShortBuffer
		}

		err = proto.Unmarshal(data, m)
		if err != nil {
			return err
		}
		return nil
	}
	return goerrors.NewInvalidNilParamError("No TypeCode() function")
}

type CreateMessageWrapperFunc func(data proto.Message) (goprotoextra.IMessageWrapper, error)

type CreateMessageFunc func() proto.Message

type TypeCodeData struct {
	CreateMessage CreateMessageFunc
	CreateWrapper CreateMessageWrapperFunc
}

var registrationMap = make(map[uint32]TypeCodeData)

func Register(tc uint32, create TypeCodeData) error {
	registrationMap[tc] = create
	return nil
}

func Find(tc uint32) (TypeCodeData, bool) {
	if tcd, ok := registrationMap[tc]; ok {
		return tcd, true
	}
	return TypeCodeData{}, false
}
