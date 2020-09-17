package test

import (
	log2 "github.com/bhbosman/gocommon/log"
	"go.uber.org/fx"
	"log"
	"testing"
)

func ProvideLogFactoryForTesting(
	t *testing.T,
	cb func(*log2.Factory)) fx.Option {
	return fx.Provide(
		func() *log2.Factory {
			logger := log.New(
				log2.NewTestWriter(t),
				"",
				log.LstdFlags)

			logFactory := log2.NewFactory(logger)
			if cb != nil {
				cb(logFactory)
			}
			return logFactory
		})
}
