package Providers

import (
	"github.com/bhbosman/gocommon/Services/implementations"
	"github.com/bhbosman/gocommon/Services/interfaces"
	"go.uber.org/fx"
)

func ProvideUniqueSessionNumber() fx.Option {
	return fx.Provide(
		func() interfaces.IUniqueSessionNumber {
			v := implementations.NewUniqueSessionNumber()
			return v
		})
}
