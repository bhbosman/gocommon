package commsImpl

import (
	"context"
	"github.com/bhbosman/gocommon/app"
	"github.com/bhbosman/gocommon/comms/common"
	"github.com/bhbosman/gocommon/comms/connectionManager"
	"go.uber.org/fx"
	"golang.org/x/sync/semaphore"
	"net"
	"net/url"
)

type netListenManager struct {
	netManager
	listener interface {
		Accept() (net.Conn, error)
	}
}

func (self *netListenManager) listenForNewConnections() {
	go func() {
		n := 0
		sem := semaphore.NewWeighted(512)
		for {
			n++
			self.logger.Printf("Trying to accept connections #%v. ", n)
			conn, err := self.Accept()
			if err != nil {
				return
			}
			if sem.TryAcquire(1) {
				self.logger.Printf("Accepted connection...")
				conn = newNetConnWithSemaphoreWrapper(conn, sem)
				self.acceptNewClientConnection(conn)
				continue
			}
			_, _ = conn.Write([]byte("ERR: To many connections\n"))
			_ = conn.Close()
		}
	}()
}

func (self *netListenManager) acceptNewClientConnection(conn net.Conn) {
	go func(conn net.Conn) {

		self.logger.Printf("Accepted %s-%s", conn.RemoteAddr(), conn.LocalAddr())
		connectionApp, ctx := self.newConnectionInstance(common.ServerConnection, conn)
		err := connectionApp.Err()
		if err != nil {
			_ = conn.Close()
			return
		}
		err = connectionApp.Start(context.Background())
		if err != nil {
			return
		}

		go func(app *fx.App, ctx context.Context) {
			<-ctx.Done()
			app.Stop(context.Background())

		}(connectionApp, ctx)
	}(conn)
}

type NewNetListenAppFunc func(params NetListenAppFuncInParams) (*fx.App, error)
type NewNetDialAppFunc func(params NetDialAppFuncInParams) (*fx.App, error)

func (self *netListenManager) Accept() (net.Conn, error) {
	return self.listener.Accept()
}

func newNetListenManager(
	params struct {
		fx.In
		Url                        *url.URL
		Listener                   net.Listener
		ConnectionReactorFactories *ConnectionReactorFactories
		ConnectionManager          connectionManager.IConnectionManager
		CancelCtx                  context.Context
		CancelFunction             context.CancelFunc
		StackFactoryFunction       TransportFactoryFunction
		Logger                     fx.ILogger
		ClientContextFactoryName   string      `name:"ConnectionReactorFactoryName"`
		ClientContext              interface{} `name:"UserContext"`
		Manager                    *app.RunTimeManager
	}) *netListenManager {
	return &netListenManager{
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
			params.ClientContext),
		listener: params.Listener,
	}
}
