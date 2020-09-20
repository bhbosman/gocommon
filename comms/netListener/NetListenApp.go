package netListener

import (
	"context"
	"fmt"
	"github.com/bhbosman/gocommon/app"
	"github.com/bhbosman/gocommon/comms/commsImpl"
	"github.com/bhbosman/gocommon/comms/connectionManager"
	"github.com/bhbosman/gocommon/log"
	"go.uber.org/fx"
	"net"
	url2 "net/url"
)

type NetListenAppFuncInParams struct {
	fx.In
	ClientContextFactories *commsImpl.ConnectionReactorFactories
	ParentContext          context.Context `name:"Application"`
	Lifecycle              fx.Lifecycle
	StackFactory           *commsImpl.TransportFactory
	Manager                *app.RunTimeManager
	ConnectionManager      connectionManager.IConnectionManager
	LogFactory             *log.Factory
}

func NewNetListenApp(
	connectionName string,
	url string,
	stackName string,
	userContextFactoryName string, settings ...ListenAppSettingsApply) NewNetListenAppFunc {
	return func(params NetListenAppFuncInParams) (*fx.App, error) {
		return fx.New(
			fx.Supply(settings),
			commsImpl.CommonComponents(
				url,
				stackName,
				params.ClientContextFactories,
				params.ParentContext,
				params.StackFactory,
				params.Manager,
				params.ConnectionManager,
				userContextFactoryName,
				params.LogFactory),

			fx.Provide(fx.Annotated{Target: newNetListenManager}),
			fx.Provide(
				func(Lifecycle fx.Lifecycle, url *url2.URL) (net.Listener, error) {
					con, err := net.Listen(url.Scheme, url.Host)
					if err != nil {
						return nil, err
					}
					Lifecycle.Append(fx.Hook{
						OnStart: nil,
						OnStop: func(ctx context.Context) error {
							return con.Close()
						},
					})
					return con, nil
				}),
			fx.Provide(
				func(params struct {
					fx.In
					Factory *log.Factory
				}) *log.SubSystemLogger {
					return params.Factory.Create(fmt.Sprintf("Listener for %v", connectionName))
				}),

			fx.Invoke(
				func(netManager *netListenManager, logger *log.SubSystemLogger, cancelFunc context.CancelFunc) {
					params.Lifecycle.Append(fx.Hook{
						OnStart: func(ctx context.Context) error {
							netManager.listenForNewConnections()
							return nil
						},
						OnStop: func(ctx context.Context) error {
							cancelFunc()
							return nil
						},
					})
				}),
		), nil
	}
}
