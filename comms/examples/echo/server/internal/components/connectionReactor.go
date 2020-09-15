package Server

import (
	"context"
	rxgo "github.com/ReactiveX/RxGo"
	"github.com/bhbosman/gocommon/comms/commsImpl"
	"github.com/bhbosman/gocommon/comms/connectionManager"
	"github.com/bhbosman/gocommon/multiBlock"
	"github.com/bhbosman/gocommon/stream"
	"go.uber.org/fx"
	"net"
	"net/url"
)

type connectionReactor struct {
	commsImpl.BaseConnectionReactor
}

func (self *connectionReactor) Init(
	conn net.Conn,
	url *url.URL,
	connectionId string,
	connectionManager connectionManager.IConnectionManager,
	onSend stream.ToConnectionFunc,
	ddd stream.ToReactorFunc) (rxgo.NextExternalFunc, error) {
	_, err := self.BaseConnectionReactor.Init(conn, url, connectionId, connectionManager, onSend, ddd)
	if err != nil {
		return nil, err
	}
	self.Logger.NameChange(self.ConnectionId)
	return self.doNext, nil
}

func (self *connectionReactor) Close() error {
	return self.BaseConnectionReactor.Close()
}

func (self *connectionReactor) Open() error {
	return self.BaseConnectionReactor.Open()
}

func (self *connectionReactor) doNext(external bool, i interface{}) {
	if self.CancelCtx.Err() != nil {
		return
	}
	defer func() {
		recover()
	}()
	switch v := i.(type) {
	case *multiBlock.ReaderWriter:
		self.ToConnection(v)
	}
}

func newConnectionReactor(
	logger fx.ILogger,
	cancelCtx context.Context,
	cancelFunc context.CancelFunc,
	name string,
	userContext interface{}) *connectionReactor {
	result := &connectionReactor{
		BaseConnectionReactor: commsImpl.NewBaseConnectionReactor(
			logger,
			name,
			cancelCtx,
			cancelFunc,
			userContext),
	}
	return result
}
