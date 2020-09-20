package netDial

import (
	"context"
	"fmt"
	"github.com/bhbosman/gocommon/app"
	"github.com/bhbosman/gocommon/comms/commsImpl"
	"github.com/bhbosman/gocommon/comms/connectionManager"
	"github.com/bhbosman/gocommon/log"
	"go.uber.org/fx"
)

type AppFuncInParams struct {
	fx.In
	ClientContextFactories *commsImpl.ConnectionReactorFactories
	ParentContext          context.Context `name:"Application"`
	Lifecycle              fx.Lifecycle
	StackFactory           *commsImpl.TransportFactory
	Manager                *app.RunTimeManager
	ConnectionManager      connectionManager.IConnectionManager
	LogFactory             *log.Factory
}
type AppFunc func(params AppFuncInParams) (*fx.App, error)

func NewNetDialApp(
	connectionName string,
	url string,
	stackName string,
	userContextFactoryName string,
	options ...DialAppSettingsApply) AppFunc {
	return func(params AppFuncInParams) (*fx.App, error) {
		return fx.New(
			fx.Supply(options),
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
			fx.Provide(
				func(params struct {
					fx.In
					Factory *log.Factory
				}) *log.SubSystemLogger {
					return params.Factory.Create(fmt.Sprintf("Dialer for %v", connectionName))
				}),
			fx.Provide(fx.Annotated{Target: newNetDialManager}),
			fx.Invoke(
				func(netManager *netDialManager, logger *log.SubSystemLogger, cancelFunction context.CancelFunc) {
					params.Lifecycle.Append(fx.Hook{
						OnStart: func(ctx context.Context) error {
							return netManager.Start(ctx)
						},
						OnStop: func(ctx context.Context) error {
							cancelFunction()
							return netManager.Stop(ctx)
						},
					})
				}),
		), nil
	}
}

