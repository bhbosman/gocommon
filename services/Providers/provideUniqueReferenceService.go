package Providers

import (
	"fmt"
	"github.com/bhbosman/gocommon/services/interfaces"
	"go.uber.org/fx"
)

func ProvideUniqueReferenceServiceInstance(UniqueSessionNumber interfaces.IUniqueReferenceService) fx.Option {
	return fx.Provide(
		fx.Annotated{
			Target: func() (interfaces.IUniqueReferenceService, error) {
				if UniqueSessionNumber == nil {
					return nil, fmt.Errorf("interfaces.IUniqueReferenceService is nil. Please resolve")
				}
				return UniqueSessionNumber, nil
			},
		},
	)
}
