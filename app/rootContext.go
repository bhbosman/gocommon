package app

import (
	"context"
	"github.com/bhbosman/gocommon"
	"github.com/bhbosman/gologging"
	"go.uber.org/fx"
)


func RegisterRootContext() fx.Option {
	return fx.Options(
		fx.Provide(
			fx.Annotated{
				Target: func(params struct {
					fx.In
					Logger        *gologging.Factory
					Lifecycle     fx.Lifecycle
					CancelContext context.Context `name:"Application"`
				}) (*gocommon.RunTimeManager, error) {
					result := gocommon.NewRunTimeManager(params.Logger, params.CancelContext)
					params.Lifecycle.Append(fx.Hook{
						OnStart: func(ctx context.Context) error {
							return result.Start(ctx)
						},
						OnStop: func(ctx context.Context) error {
							return result.Stop(ctx)
						},
					})
					return result, nil
				}}),
		fx.Provide(
			fx.Annotated{
				Name: "Application",
				Target: func(params struct {
					fx.In
					Lifecycle fx.Lifecycle
				}) (context.Context, context.CancelFunc) {
					cancel, cancelFunc := context.WithCancel(context.Background())
					params.Lifecycle.Append(fx.Hook{
						OnStart: nil,
						OnStop: func(ctx context.Context) error {
							cancelFunc()
							return nil
						},
					})
					return cancel, cancelFunc
				}}),
		ProvidePubSub("Application"))
}
