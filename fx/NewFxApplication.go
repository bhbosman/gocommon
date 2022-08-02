package fx

import (
	"go.uber.org/fx"
	"time"
)

func NewFxApplicationOptions(
	startTimeOut time.Duration,
	stopTimeOut time.Duration,
	option ...fx.Option,
) fx.Option {
	return fx.Options(
		ProvideFxWithLogger(),
		fx.StartTimeout(startTimeOut),
		fx.StopTimeout(stopTimeOut),
		fx.Options(option...),
	)
}
