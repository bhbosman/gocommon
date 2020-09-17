package commsImpl

import (
	"context"
	"github.com/bhbosman/gocommon/app"
	"github.com/bhbosman/gocommon/comms/common"
	"github.com/bhbosman/gocommon/comms/connectionManager"
	"github.com/bhbosman/gocommon/log"
	"go.uber.org/fx"
	"golang.org/x/sync/semaphore"
	"net"
	"net/url"
	"time"
)

type netDialManager struct {
	netManager
}

func (self *netDialManager) Start(_ context.Context) error {
	go func() {
		sem := semaphore.NewWeighted(1)
		for {
			if self.cancelCtx.Err() != nil {
				return
			}
			if sem.Acquire(self.cancelCtx, 1) != nil {
				return
			}
			dialer := net.Dialer{
				Timeout: time.Second * 30,
			}
			conn, err := dialer.DialContext(self.cancelCtx, "tcp4", self.url.Host)
			if err != nil {
				sem.Release(1)
				continue
			}
			if self.cancelCtx.Err() != nil {
				_ = conn.Close()
				return
			}
			conn = newNetConnWithSemaphoreWrapper(conn, sem)
			instance, ctx := self.newConnectionInstance(common.ClientConnection, conn)
			if instance.Err() != nil {
				_ = conn.Close()
				continue
			}
			err = instance.Start(context.Background())
			if err != nil {
				continue
			}
			go func(app *fx.App, ctx context.Context) {
				<-ctx.Done()
				_ = app.Stop(context.Background())
			}(instance, ctx)
		}
	}()
	return nil
}

func (self *netDialManager) Stop(_ context.Context) error {
	return nil
}

func newNetDialManager(
	params struct {
		fx.In
		Url                        *url.URL
		ConnectionReactorFactories *ConnectionReactorFactories
		ConnectionManager          connectionManager.IConnectionManager
		CancelCtx                  context.Context
		CancelFunction             context.CancelFunc
		StackFactoryFunction       TransportFactoryFunction
		Logger                     *log.SubSystemLogger
		ClientContextFactoryName   string      `name:"ConnectionReactorFactoryName"`
		ClientContext              interface{} `name:"UserContext"`
		Manager                    *app.RunTimeManager
		LogFactory                 *log.Factory
	}) *netDialManager {
	return &netDialManager{
		netManager: newNetManager(
			params.Url,
			params.ConnectionReactorFactories,
			params.CancelCtx,
			params.CancelFunction,
			params.Logger,
			params.StackFactoryFunction,
			params.ClientContextFactoryName,
			params.Manager,
			params.ConnectionManager,
			params.LogFactory,
			params.ClientContext),
	}
}
