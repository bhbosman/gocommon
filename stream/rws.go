package stream

import (
	"context"
	"encoding/binary"
	"errors"
	"github.com/bhbosman/gocommon/constants"
	"github.com/bhbosman/gocommon/multiBlock"
	"google.golang.org/protobuf/proto"
	"io"
)

type IReadWriterSize interface {
	io.ReadWriter
	Size() int
}

type IReaderSize interface {
	io.Reader
	Size() int
}

type ReadWriterSize IReadWriterSize
type ReaderSize IReaderSize

type IMessageWrapper interface {
	CancelCtx() context.Context
	CancelFunc() context.CancelFunc
	ToReactor(inline bool, any interface{}) error
	ToConnection(rws ReadWriterSize) error
	Message() proto.Message
}
type BaseMessageWrapper struct {
	cancelCtx    context.Context
	cancelFunc   context.CancelFunc
	toReactor    ToReactorFunc
	toConnection func(rws ReadWriterSize) error
}
type ToReactorFunc func(inline bool, any interface{}) error
type ToConnectionFunc func(rws ReadWriterSize) error
type ConnectionReactorHandler func(i interface{})
type ErrorStateFunc func() (bool, error)

func NewBaseMessageWrapper(
	cancelCtx context.Context,
	cancelFunc context.CancelFunc,
	toReactor ToReactorFunc,
	toConnection ToConnectionFunc) BaseMessageWrapper {
	return BaseMessageWrapper{
		cancelCtx:    cancelCtx,
		cancelFunc:   cancelFunc,
		toReactor:    toReactor,
		toConnection: toConnection,
	}
}

var NullValueError = errors.New("NullValueError")

func (self *BaseMessageWrapper) ToReactor(inline bool, any interface{}) error {
	if self != nil && self.toReactor != nil {
		if self.cancelCtx != nil {
			err := self.cancelCtx.Err()
			if err != nil {
				return err
			}
		}
		return self.toReactor(inline, any)
	}
	return NullValueError
}

func (self *BaseMessageWrapper) ToConnection(rws ReadWriterSize) error {
	if self != nil && self.toConnection != nil {
		if self.cancelCtx != nil {
			err := self.cancelCtx.Err()
			if err != nil {
				return err
			}
		}
		return self.toConnection(rws)
	}
	return NullValueError
}

func (self *BaseMessageWrapper) CancelCtx() context.Context {
	if self != nil && self.cancelCtx != nil {
		return self.cancelCtx
	}
	return nil
}

func (self *BaseMessageWrapper) CancelFunc() context.CancelFunc {
	if self != nil && self.cancelFunc != nil {
		return self.cancelFunc
	}
	return nil
}

type MessageWrapper struct {
	BaseMessageWrapper
	Data proto.Message
}

//func ToConnection(description string, logger fx.ILogger, cancelCtx context.Context, toConnection ToConnectionFunc) ToConnectionFunc {
//	var exitError error
//	logError := func() {
//		logger.Error(description, exitError)
//	}
//	return func(rws ReadWriterSize) error {
//		if exitError != nil {
//			return exitError
//		}
//		exitError = cancelCtx.Err()
//		if exitError != nil {
//			logError()
//			return exitError
//		}
//		exitError = toConnection(rws)
//		if exitError != nil {
//			logError()
//			return exitError
//		}
//		return exitError
//	}
//}
//func ToReactorWithContext(description string, logger fx.ILogger, cancelCtx context.Context, toReactor ToReactorFunc) ToReactorFunc {
//	var exitError error
//	logError := func() {
//		logger.Error(description, exitError)
//	}
//	return func(i interface{}) error {
//		if exitError != nil {
//			return exitError
//		}
//		exitError = cancelCtx.Err()
//		if exitError != nil {
//			logError()
//			return exitError
//		}
//		exitError = toReactor(i)
//		if exitError != nil {
//			logError()
//			return exitError
//		}
//		return exitError
//	}
//}

func NewMessageWrapper(
	cancelCtx context.Context,
	cancelFunc context.CancelFunc,
	toReactor ToReactorFunc,
	toConnection ToConnectionFunc,
	data proto.Message) *MessageWrapper {
	return &MessageWrapper{
		BaseMessageWrapper: NewBaseMessageWrapper(
			cancelCtx,
			cancelFunc,
			toReactor,
			toConnection),
		Data: data,
	}
}

func (self *MessageWrapper) Message() proto.Message {
	return self.Data
}

func (self *MessageWrapper) messageWrapper() IMessageWrapper {
	return self
}

func Marshall(m proto.Message) (stm ReadWriterSize, err error) {
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
	toReactor ToReactorFunc,
	toConnection ToConnectionFunc) (msgWrapper IMessageWrapper, err error) {
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

func UnMarshalMessage(rws ReadWriterSize, m proto.Message) error {
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
	toReactor ToReactorFunc,
	toConnection ToConnectionFunc,
	data proto.Message) (IMessageWrapper, error)

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
