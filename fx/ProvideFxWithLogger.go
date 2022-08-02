package fx

import (
	"go.uber.org/fx"
	"go.uber.org/fx/fxevent"
	"go.uber.org/zap"
)

func ProvideFxWithLogger() fx.Option {
	return fx.WithLogger(
		func(params struct {
			fx.In
			Logger *zap.Logger
		}) fxevent.Logger {
			return &fxevent.ZapLogger{
				Logger: params.Logger.Named("FX"),
			}
		})
}
