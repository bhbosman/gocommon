package netDial

import (
	"context"
	"github.com/bhbosman/gocommon/app"
	"github.com/bhbosman/gocommon/comms/common"
	"github.com/bhbosman/gocommon/comms/commsImpl"
	"github.com/bhbosman/gocommon/comms/connectionManager"
	"github.com/bhbosman/gocommon/log"
	"go.uber.org/fx"
	"golang.org/x/sync/semaphore"
	"net"
	"net/url"
	"time"
)

type ICanDial interface {
	CanDial() bool
	DialSuccess()
}

type netDialManager struct {
	commsImpl.NetManager
	CanDial []ICanDial
}

func (self *netDialManager) Start(_ context.Context) error {
	go func() {
		sem := semaphore.NewWeighted(1)
	loop:
		for {
			if len(self.CanDial) > 0 {
				b := true
				for _, canDial := range self.CanDial {
					b = b && canDial.CanDial()
					if !b {
						time.Sleep(time.Second)
						continue loop
					}
				}
			}

			if self.CancelCtx.Err() != nil {
				return
			}
			if sem.Acquire(self.CancelCtx, 1) != nil {
				return
			}
			dialer := net.Dialer{
				Timeout: time.Second * 30,
			}
			conn, err := dialer.DialContext(self.CancelCtx, "tcp4", self.Url.Host)
			if err != nil {
				sem.Release(1)
				continue loop
			}
			if self.CancelCtx.Err() != nil {
				_ = conn.Close()
				return
			}
			conn = commsImpl.NewNetConnWithSemaphoreWrapper(conn, sem)
			instance, ctx := self.NewConnectionInstance(common.ClientConnection, conn)
			if instance.Err() != nil {
				_ = conn.Close()
				continue loop
			}
			err = instance.Start(context.Background())
			if err != nil {
				continue loop
			}
			go func(app *fx.App, ctx context.Context) {
				ticker := time.NewTicker(time.Second)
				defer ticker.Stop()
				for {
					select {
					case <-ctx.Done():
						_ = app.Stop(context.Background())
						return
					case <-ticker.C:
						b := true
						for _, canDial := range self.CanDial {
							b = b && canDial.CanDial()
							if !b {
								_ = app.Stop(context.Background())
								return
							}
						}
					}
				}

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
		ConnectionReactorFactories *commsImpl.ConnectionReactorFactories
		ConnectionManager          connectionManager.IConnectionManager
		CancelCtx                  context.Context
		CancelFunction             context.CancelFunc
		StackFactoryFunction       commsImpl.TransportFactoryFunction
		Logger                     *log.SubSystemLogger
		ClientContextFactoryName   string `name:"ConnectionReactorFactoryName"`
		Manager                    *app.RunTimeManager
		LogFactory                 *log.Factory
		Options                    []DialAppSettingsApply
	}) *netDialManager {

	settings := &dialAppSettings{
		userContext: nil,
		canDial:     nil,
	}
	for _, option := range params.Options {
		option.apply(settings)
	}

	return &netDialManager{
		NetManager: commsImpl.NewNetManager(
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
			settings.userContext),
		CanDial: settings.canDial,
	}
}
