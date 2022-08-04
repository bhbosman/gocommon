package model

import (
	"context"
	"github.com/icza/gox/fmtx"
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
	Id               string
	CancelContext    context.Context
	CancelFunc       context.CancelFunc
	Name             string
	ConnectionTime   time.Time
	InboundCounters  *PublishRxHandlerCounters
	OutboundCounters *PublishRxHandlerCounters
	Grid             []*LineData
	KeyValuesMap     map[string]string
}

func NewConnectionInformation(id string, function context.CancelFunc, CancelContext context.Context) *ConnectionInformation {
	return &ConnectionInformation{
		Id:             id,
		CancelContext:  CancelContext,
		CancelFunc:     function,
		ConnectionTime: time.Now().Truncate(time.Second),
		//StackProperties: make(map[StackPropertyKey]StackPropertyValue),
	}
}

type ConnectionCreated struct {
	ConnectionId   string
	ConnectionName string
	CancelFunc     context.CancelFunc
	ConnectionTime time.Time
	CancelContext  context.Context
}

func NewConnectionCreated(
	connectionId string,
	connectionName string,
	cancelFunc context.CancelFunc,
	connectionTime time.Time,
	cancelContext context.Context,
) *ConnectionCreated {
	return &ConnectionCreated{
		ConnectionId:   connectionId,
		ConnectionName: connectionName,
		CancelFunc:     cancelFunc,
		ConnectionTime: connectionTime,
		CancelContext:  cancelContext,
	}
}

type ConnectionClosed struct {
	ConnectionId string
}

type ConnectionState struct {
	ConnectionId   string
	CancelContext  context.Context
	CancelFunc     context.CancelFunc
	Name           string
	ConnectionTime time.Time
	Grid           []LineData
	KeyValue       []KeyValue
}
