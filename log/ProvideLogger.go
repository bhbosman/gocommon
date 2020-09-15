package log

import (
	"go.uber.org/fx"
	"log"
)

func ProvideLogger(logger *log.Logger) fx.Option {
	return fx.Provide(
		func() *log.Logger {
			return logger
		})
}
