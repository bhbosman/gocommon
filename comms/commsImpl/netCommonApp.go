package commsImpl

import (
	"context"
	rxgo "github.com/ReactiveX/RxGo"
	"github.com/bhbosman/gocommon/app"
	"github.com/bhbosman/gocommon/comms/common"
	"github.com/bhbosman/gocommon/comms/connectionManager"
	"github.com/bhbosman/gocommon/log"
	"go.uber.org/fx"
)

func CommonComponents(
	url string,
	stackName string,
	ClientContextFactories *ConnectionReactorFactories,
	ParentContext context.Context,
	StackFactory *TransportFactory,
	Manager *app.RunTimeManager,
	ConnectionManager connectionManager.IConnectionManager,
	connectionReactorFactoryName string,
	logFactory *log.Factory,
	userContext interface{}) fx.Option {
	return fx.Options(
		fx.Supply(ClientContextFactories, StackFactory, Manager, logFactory),
		fx.Provide(fx.Annotated{Name: UserContext, Target: func() interface{} {
			return userContext
		}}),
		fx.Provide(fx.Annotated{Target: common.CreateUrl(url)}),
		fx.Provide(fx.Annotated{Name: "StackName", Target: CreateStringContext(stackName)}),
		fx.Provide(fx.Annotated{Name: ConnectionReactorFactoryName, Target: CreateStringContext(connectionReactorFactoryName)}),
		fx.Provide(fx.Annotated{Target: func() connectionManager.IConnectionManager { return ConnectionManager }}),
		fx.Provide(fx.Annotated{Target: func() connectionManager.IRegisterToConnectionManager { return ConnectionManager }}),
		fx.Provide(fx.Annotated{Target: func() connectionManager.IObtainConnectionManagerInformation { return ConnectionManager }}),
		fx.Provide(fx.Annotated{Target: func() connectionManager.ICommandsToConnectionManager { return ConnectionManager }}),
		fx.Provide(fx.Annotated{Target: func() rxgo.IPublishToConnectionManager { return ConnectionManager }}),
		fx.Provide(fx.Annotated{Target: func() (ctx context.Context, cancel context.CancelFunc) {
			return context.WithCancel(ParentContext)
		}}),
		fx.Provide(fx.Annotated{
			Target: func(
				params struct {
					fx.In
					StackName string `name:"StackName"`
					Factory   *TransportFactory
				}) (TransportFactoryFunction, error) {
				return params.Factory.Get(params.StackName)
			}}),
	)
}
