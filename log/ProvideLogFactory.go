package log

import (
	"go.uber.org/fx"
	"log"
)

func ProvideLogFactory(logger *log.Logger, cb func(*Factory)) fx.Option {
	return fx.Provide(
		func() *Factory {
			logFactory := NewFactory(logger)
			if cb != nil {
				cb(logFactory)
			}
			return logFactory
		})
}
