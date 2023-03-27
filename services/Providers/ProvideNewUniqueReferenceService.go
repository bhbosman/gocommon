package Providers

import (
	"github.com/bhbosman/gocommon/services/implementations"
	"github.com/bhbosman/gocommon/services/interfaces"
	"go.uber.org/fx"
)

func ProvideNewUniqueReferenceService() fx.Option {
	return fx.Provide(
		func(uniqueSessionNumber interfaces.IUniqueSessionNumber) interfaces.IUniqueReferenceService {
			v := implementations.NewUniqueReferenceService(uniqueSessionNumber)
			return v
		})
}
