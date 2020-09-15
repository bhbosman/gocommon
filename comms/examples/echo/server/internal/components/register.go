package Server

import (
	"github.com/bhbosman/gocommon/comms/commsImpl"
	"go.uber.org/fx"
)

func RegisterEchoServiceListener() fx.Option {
	const CreateClientHandlerFactoryName = "EchoServerConnectionReactorFactory"
	return fx.Options(
		fx.Provide(
			fx.Annotated{
				Group: commsImpl.ConnectionReactorFactoryConst,
				Target: func() (commsImpl.IConnectionReactorFactory, error) {
					return &connectionReactorFactory{
						name: CreateClientHandlerFactoryName,
					}, nil
				},
			}),
		fx.Provide(
			fx.Annotated{
				Group: "Apps",
				Target: commsImpl.NewNetListenApp(
					"EchoServerConnectionManager(Empty)",
					"tcp4://127.0.0.1:3000",
					commsImpl.TransportFactoryEmptyName,
					CreateClientHandlerFactoryName,
					nil),
			}),
		fx.Provide(
			fx.Annotated{
				Group: "Apps",
				Target: commsImpl.NewNetListenApp(
					"EchoServerConnectionManager(Compressed)",
					"tcp4://127.0.0.1:3001",
					commsImpl.TransportFactoryCompressedName,
					CreateClientHandlerFactoryName,
					nil),
			}),
		fx.Provide(
			fx.Annotated{
				Group: "Apps",
				Target: commsImpl.NewNetListenApp(
					"EchoServerConnectionManager(UnCompressed)",
					"tcp4://127.0.0.1:3002",
					commsImpl.TransportFactoryUnCompressedName,
					CreateClientHandlerFactoryName,
					nil),
			}),
	)
}
