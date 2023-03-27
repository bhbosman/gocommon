package Providers

import (
	"github.com/bhbosman/gocommon/services/implementations"
	"github.com/bhbosman/gocommon/services/interfaces"
	"go.uber.org/fx"
)

func ProvideUniqueSessionNumber() fx.Option {
	return fx.Provide(
		func() interfaces.IUniqueSessionNumber {
			v := implementations.NewUniqueSessionNumber()
			return v
		})
}
