package log

import (
	"go.uber.org/fx"
	"log"
)

func ProvideLogFactory(logger *log.Logger, cb func(*LogFactory)) fx.Option {
	return fx.Provide(
		func() *LogFactory {
			logFactory := NewLogFactory(logger)
			if cb != nil {
				cb(logFactory)
			}
			return logFactory
		})
}
