package GoFunctionCounter

import (
	"context"
	"github.com/bhbosman/gocommon/Services/interfaces"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

func Provide() fx.Option {
	return fx.Options(
		fx.Provide(
			fx.Annotated{
				Target: func(
					params struct {
						fx.In
					},
				) (func() (IData, error), error) {
					return func() (IData, error) {
						return newData()
					}, nil
				},
			},
		),
		fx.Provide(
			fx.Annotated{
				Target: func(
					params struct {
						fx.In
						ApplicationContext     context.Context `name:"Application"`
						Lifecycle              fx.Lifecycle
						OnData                 func() (IData, error)
						Logger                 *zap.Logger
						UniqueReferenceService interfaces.IUniqueReferenceService
						UniqueSessionNumber    interfaces.IUniqueSessionNumber
					},
				) (IService, error) {
					serviceInstance, err := newService(
						params.ApplicationContext,
						params.OnData,
						params.Logger,
						params.UniqueReferenceService,
						params.UniqueSessionNumber,
					)
					if err != nil {
						return nil, err
					}
					params.Lifecycle.Append(
						fx.Hook{
							OnStart: serviceInstance.OnStart,
							OnStop:  serviceInstance.OnStop,
						})
					return serviceInstance, nil
				},
			},
		),
	)
}
