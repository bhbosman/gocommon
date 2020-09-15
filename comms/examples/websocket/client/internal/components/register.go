package components

import (
	"github.com/bhbosman/goConnect/polygon/stream"
	"github.com/bhbosman/gocommon/comms/commsImpl"
	"go.uber.org/fx"
)

func RegisterWebSocketDialerForPolygonForex() fx.Option {
	const createServerHandlerFactoryName = "ForexConnectionReactorFactory"
	return fx.Options(
		fx.Provide(fx.Annotated{
			Group: commsImpl.ConnectionReactorFactoryConst,
			Target: func(
				params struct {
					fx.In
					ApiKey string `name:"PolygonApiKey"`
				}) (commsImpl.IConnectionReactorFactory, error) {

				return stream.NewConnectionReactorForexFactory(
					createServerHandlerFactoryName,
					params.ApiKey), nil
			},
		}),
		fx.Provide(fx.Annotated{
			Group: "Apps",
			Target: commsImpl.NewNetDialApp(
				"socket.polygon.io:443",
				"wss://socket.polygon.io:443/forex",
				commsImpl.WebSocketName,
				createServerHandlerFactoryName,
				nil),
		}),
	)
}
