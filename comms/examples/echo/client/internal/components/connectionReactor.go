package components

import (
	"context"
	"fmt"
	rxgo "github.com/ReactiveX/RxGo"
	"github.com/bhbosman/gocommon/comms/commsImpl"
	"github.com/bhbosman/gocommon/comms/connectionManager"
	"github.com/bhbosman/gocommon/multiBlock"
	"github.com/bhbosman/gocommon/stream"
	"go.uber.org/fx"
	"net"
	"net/url"

	"time"
)

type DoPublish struct {
}

type connectionReactor struct {
	commsImpl.BaseConnectionReactor
	publish         *time.Ticker
	sendPackets     int
	receivePackets  int
	sendPacketSz    int
	receivePacketSz int
}

func (self *connectionReactor) Init(
	conn net.Conn,
	url *url.URL,
	connectionId string,
	connectionManager connectionManager.IConnectionManager,
	onSend stream.ToConnectionFunc,
	toConnectionReactor stream.ToReactorFunc) (rxgo.NextExternalFunc, error) {
	_, _ = self.BaseConnectionReactor.Init(conn, url, connectionId, connectionManager, onSend, toConnectionReactor)
	self.Logger.NameChange(self.ConnectionId)
	return self.doNext, nil

}

func (self *connectionReactor) Close() error {
	err := self.BaseConnectionReactor.Close()
	if err != nil {
		return err
	}

	if self.publish != nil {
		self.publish.Stop()
	}
	return nil
}

func (self *connectionReactor) Open() error {
	err := self.BaseConnectionReactor.Open()
	if err != nil {
		return err
	}

	ticker := time.NewTimer(time.Second)
	self.publish = time.NewTicker(time.Second * 10)
	go func() {
		defer func() {
			r := recover()
			if r != nil {
				print(r)
			}

		}()
		for {
			select {
			case <-self.CancelCtx.Done():
				return
			case _, ok := <-ticker.C:
				if !ok {
					return
				}
				ticker.Stop()
				if self.CancelCtx.Err() != nil {
					return
				}
				data := make([]byte, 1024*1024)
				err := self.ToConnection(multiBlock.NewReaderWriterBlock(data))
				if err != nil {
					return
				}
				self.sendPackets++
				self.sendPacketSz += len(data)
				ticker.Reset(time.Millisecond * 8)
			case _, ok := <-self.publish.C:
				if !ok {
					return
				}
				if self.CancelCtx.Err() != nil {
					return
				}
				self.ToReactor(false, &DoPublish{})
			}
		}
	}()
	return nil
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

func (self *connectionReactor) doNext(external bool, i interface{}) {
	switch v := i.(type) {
	case *DoPublish:
		self.Logger.Printf(fmt.Sprintf("%10d, %15d, %10d, %15d, %15d", self.sendPackets, self.sendPacketSz, self.receivePackets, self.receivePacketSz, self.sendPackets-self.receivePackets))
	case *multiBlock.ReaderWriter:
		self.receivePackets++
		self.receivePacketSz += v.Size()
	default:
	}
}
