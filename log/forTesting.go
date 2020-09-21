package log

import (
	"go.uber.org/fx"
	"log"
	"testing"
)

func ProvideLogFactoryForTesting(
	t *testing.T,
	cb func(*Factory)) fx.Option {
	return fx.Provide(
		func() *Factory {
			logger := log.New(
				NewTestWriter(t),
				"",
				log.LstdFlags)

			logFactory := NewFactory(logger)
			if cb != nil {
				cb(logFactory)
			}
			return logFactory
		})
}
