package Providers

import (
	"context"
	"github.com/bhbosman/gocommon/GoFunctionCounter"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

func RegisterRunTimeManager() fx.Option {
	return fx.Provide(
		fx.Annotated{
			Target: func(params struct {
				fx.In
				Logger            *zap.Logger
				Lifecycle         fx.Lifecycle
				GoFunctionCounter GoFunctionCounter.IService
				CancelContext     context.Context `name:"Application"`
			}) (*RunTimeManager, error) {
				result := NewRunTimeManager(
					params.Logger,
					params.CancelContext,
					params.GoFunctionCounter,
				)
				params.Lifecycle.Append(fx.Hook{
					OnStart: func(ctx context.Context) error {
						return result.Start(ctx)
					},
					OnStop: func(ctx context.Context) error {
						return result.Stop(ctx)
					},
				})
				return result, nil
			}})
}
