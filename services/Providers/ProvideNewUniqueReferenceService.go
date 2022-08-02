package Providers

import (
	"github.com/bhbosman/gocommon/Services/implementations"
	"github.com/bhbosman/gocommon/Services/interfaces"
	"go.uber.org/fx"
)

func ProvideNewUniqueReferenceService() fx.Option {
	return fx.Provide(
		func(uniqueSessionNumber interfaces.IUniqueSessionNumber) interfaces.IUniqueReferenceService {
			v := implementations.NewUniqueReferenceService(uniqueSessionNumber)
			return v
		})
}
