package stream

import (
	"context"
	"encoding/binary"
	"github.com/bhbosman/gocommon/constants"
	"github.com/bhbosman/gocommon/multiBlock"
	"github.com/bhbosman/goprotoextra"
	"google.golang.org/protobuf/proto"
)

func Marshall(m proto.Message) (stm goprotoextra.ReadWriterSize, err error) {
	if tc, ok := m.(interface {
		TypeCode() uint32
	}); ok {
		tcBytes := [8]byte{}
		binary.LittleEndian.PutUint32(tcBytes[0:4], tc.TypeCode())
		binary.LittleEndian.PutUint32(tcBytes[4:8], uint32(proto.Size(m)))
		var marshallBytes []byte
		marshallBytes, err = proto.Marshal(m)
		if err != nil {
			return nil, err
		}
		stm = multiBlock.NewReaderWriterWithBlocks(tcBytes[:], marshallBytes)
		return stm, err
	}
	return nil, constants.InvalidParam
}

func UnMarshal(
	rws *multiBlock.ReaderWriter,
	cancelCtx context.Context,
	cancelFunc context.CancelFunc,
	toReactor goprotoextra.ToReactorFunc,
	toConnection goprotoextra.ToConnectionFunc) (msgWrapper goprotoextra.IMessageWrapper, err error) {
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
		return r.CreateWrapper(cancelCtx, cancelFunc, toReactor, toConnection, msg)
	}
	return nil, constants.InvalidParam
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
			return constants.InvalidParam
		}
		n := binary.LittleEndian.Uint32(tcBytes[4:8])
		if n < uint32(rws.Size()) {
			return constants.InvalidParam
		}
		data := make([]byte, n)
		nn, err := rws.Read(data)
		if err != nil {
			return err
		}
		if uint32(nn) != n {
			return constants.InvalidParam
		}

		err = proto.Unmarshal(data, m)
		if err != nil {
			return err
		}
		return nil
	}
	return constants.InvalidParam
}

type CreateMessageWrapperFunc func(
	cancelCtx context.Context,
	cancelFunc context.CancelFunc,
	toReactor goprotoextra.ToReactorFunc,
	toConnection goprotoextra.ToConnectionFunc,
	data proto.Message) (goprotoextra.IMessageWrapper, error)

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
