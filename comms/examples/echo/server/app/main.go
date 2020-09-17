package main

import (
	"context"
	app2 "github.com/bhbosman/gocommon/app"
	"github.com/bhbosman/gocommon/comms/commsImpl"
	"github.com/bhbosman/gocommon/comms/connectionManager"
	"github.com/bhbosman/gocommon/comms/connectionManager/endpoints"
	"github.com/bhbosman/gocommon/comms/connectionManager/view"
	echoServer "github.com/bhbosman/gocommon/comms/examples/echo/server/internal/components"
	"github.com/bhbosman/gocommon/comms/http"
	log2 "github.com/bhbosman/gocommon/log"
	"go.uber.org/fx"
	"log"
	"os"
)

func main() {
	app := fx.New(
		log2.ProvideLogFactory(log.New(os.Stderr, "EchoServer: ", log.LstdFlags), nil),
		connectionManager.RegisterDefaultConnectionManager(),
		commsImpl.RegisterAllConnectionRelatedServices(),
		endpoints.RegisterConnectionManagerEndpoint(),
		view.RegisterConnectionsHtmlTemplate(),
		http.RegisterHttpHandler("http://127.0.0.1:8080"),
		app2.RegisterRootContext(),
		echoServer.RegisterEchoServiceListener(),
		fx.Provide(
			func(params struct {
				fx.In
				Factory *log2.Factory
			}) *log2.SubSystemLogger {
				return params.Factory.Create("Main")
			}),

		fx.Invoke(
			func(params struct {
				fx.In
				Lifecycle      fx.Lifecycle
				Apps           []*fx.App `group:"Apps"`
				Logger         *log2.SubSystemLogger
				RunTimeManager *app2.RunTimeManager
			}) {
				params.Lifecycle.Append(fx.Hook{
					OnStart: func(ctx context.Context) error {
						return params.RunTimeManager.Start(ctx)
					},
					OnStop: func(ctx context.Context) error {
						return params.RunTimeManager.Stop(ctx)
					},
				})
				for _, item := range params.Apps {
					localApp := item
					params.Lifecycle.Append(fx.Hook{
						OnStart: func(ctx context.Context) error {
							return localApp.Start(ctx)
						},
						OnStop: func(ctx context.Context) error {
							return localApp.Stop(ctx)
						},
					})
				}
			}),
	)
	if app.Err() != nil {
		return
	}
	app.Run()
}
