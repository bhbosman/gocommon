package commsImpl

import (
	"context"
	"fmt"
	"github.com/bhbosman/gocommon/app"
	"github.com/bhbosman/gocommon/comms/connectionManager"
	"go.uber.org/fx"
	"time"
)

type NetDialAppFuncInParams struct {
	fx.In
	ClientContextFactories *ConnectionReactorFactories
	ParentContext          context.Context `name:"Application"`
	Lifecycle              fx.Lifecycle
	StackFactory           *TransportFactory
	Manager                *app.RunTimeManager
	ConnectionManager      connectionManager.IConnectionManager
}

func NewNetDialApp(
	connectionName string,
	url string,
	stackName string,
	userContextFactoryName string,
	userContext interface{}) NewNetDialAppFunc {
	return func(params NetDialAppFuncInParams) (*fx.App, error) {
		return fx.New(
			fx.StopTimeout(time.Hour),
			fx.LogName(fmt.Sprintf("%v", connectionName)),
			CommonComponents(
				url,
				stackName,
				params.ClientContextFactories,
				params.ParentContext,
				params.StackFactory,
				params.Manager,
				params.ConnectionManager,
				userContextFactoryName,
				userContext),
			fx.Provide(fx.Annotated{Target: newNetDialManager}),
			fx.Invoke(
				func(netManager *netDialManager, logger fx.ILogger, cancelFunction context.CancelFunc) {
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
