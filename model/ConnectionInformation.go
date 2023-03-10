package model

import (
	"context"
	"github.com/icza/gox/fmtx"
	"github.com/reactivex/rxgo/v2"
	"time"
)

type StackPropertyKey struct {
	Index     int
	Name      string
	Direction StreamDirection
}

type StackPropertyValue struct {
	msgCount  int
	byteCount int
}

func (self StackPropertyValue) MsgCount() string {
	return fmtx.FormatInt(int64(self.msgCount), 3, ',')
}

func (self StackPropertyValue) ByteCount() string {
	return fmtx.FormatSize(int64(self.byteCount), fmtx.SizeUnitAuto, 2)
}

type ConnectionInformation struct {
	Id                      string
	CancelContext           context.Context
	CancelFunc              context.CancelFunc
	Name                    string
	ConnectionTime          time.Time
	NextFuncOutBoundChannel rxgo.NextFunc
	NextFuncInBoundChannel  rxgo.NextFunc
	InboundCounters         *PublishRxHandlerCounters
	OutboundCounters        *PublishRxHandlerCounters
	Grid                    []*LineData
	KeyValuesMap            map[string]string
}

func NewConnectionInformation(
	id string,
	function context.CancelFunc,
	CancelContext context.Context,
	nextFuncOutBoundChannel rxgo.NextFunc,
	nextFuncInBoundChannel rxgo.NextFunc,

) *ConnectionInformation {
	return &ConnectionInformation{
		Id:                      id,
		CancelContext:           CancelContext,
		CancelFunc:              function,
		ConnectionTime:          time.Now().Truncate(time.Second),
		NextFuncOutBoundChannel: nextFuncOutBoundChannel,
		NextFuncInBoundChannel:  nextFuncInBoundChannel,
	}
}

type ConnectionCreated struct {
	ConnectionId            string
	ConnectionName          string
	CancelFunc              context.CancelFunc
	ConnectionTime          time.Time
	CancelContext           context.Context
	nextFuncOutBoundChannel rxgo.NextFunc
	nextFuncInBoundChannel  rxgo.NextFunc
}

func NewConnectionCreated(
	connectionId string,
	connectionName string,
	cancelFunc context.CancelFunc,
	connectionTime time.Time,
	cancelContext context.Context,
	nextFuncOutBoundChannel rxgo.NextFunc,
	nextFuncInBoundChannel rxgo.NextFunc,
) *ConnectionCreated {
	return &ConnectionCreated{
		ConnectionId:            connectionId,
		ConnectionName:          connectionName,
		CancelFunc:              cancelFunc,
		ConnectionTime:          connectionTime,
		CancelContext:           cancelContext,
		nextFuncOutBoundChannel: nextFuncOutBoundChannel,
		nextFuncInBoundChannel:  nextFuncInBoundChannel,
	}
}

type ConnectionClosed struct {
	ConnectionId string
}

type ConnectionState struct {
	ConnectionId string
	//CancelContext  context.Context
	//CancelFunc     context.CancelFunc
	//Name           string
	//ConnectionTime time.Time
	Grid     []LineData
	KeyValue []KeyValue
}
