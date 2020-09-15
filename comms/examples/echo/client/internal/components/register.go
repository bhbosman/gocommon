package components

import (
	"github.com/bhbosman/gocommon/comms/commsImpl"
	"go.uber.org/fx"
)

func RegisterEchoServiceDialer() fx.Option {
	const createServerHandlerFactoryName = "EchoClientConnectionReactorFactory"
	return fx.Options(
		fx.Provide(fx.Annotated{
			Group: commsImpl.ConnectionReactorFactoryConst,
			Target: func() (commsImpl.IConnectionReactorFactory, error) {
				return &connectionReactorFactory{
					name: createServerHandlerFactoryName,
				}, nil

			},
		}),
		fx.Provide(fx.Annotated{
			Group: "Apps",
			Target: commsImpl.NewNetDialApp(
				"EchoServiceDialer(Empty)",
				"tcp4://127.0.0.1:3000",
				commsImpl.TransportFactoryEmptyName,
				createServerHandlerFactoryName,
				nil),
		}),
		fx.Provide(fx.Annotated{
			Group: "Apps",
			Target: commsImpl.NewNetDialApp(
				"EchoServiceDialer(Compressed)",
				"tcp4://127.0.0.1:3001",
				commsImpl.TransportFactoryCompressedName,
				createServerHandlerFactoryName,
				nil),
		}),
		fx.Provide(fx.Annotated{
			Group: "Apps",
			Target: commsImpl.NewNetDialApp(
				"EchoServiceDialer(UnCompressed)",
				"tcp4://127.0.0.1:3002",
				commsImpl.TransportFactoryUnCompressedName,
				createServerHandlerFactoryName, nil),
		}),
	)
}
