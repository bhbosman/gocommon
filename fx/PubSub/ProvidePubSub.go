package PubSub

import (
	"context"
	"github.com/cskr/pubsub"
	"go.uber.org/fx"
)

func ProvidePubSub(name string) fx.Option {
	return fx.Provide(
		fx.Annotated{
			Name: name,
			Target: func(
				params struct {
					fx.In
					Lifecycle fx.Lifecycle
				},
			) *pubsub.PubSub {
				result := pubsub.New(32)
				params.Lifecycle.Append(
					fx.Hook{
						OnStart: nil,
						OnStop: func(ctx context.Context) error {
							result.Shutdown()
							return nil
						},
					},
				)
				return result
			},
		},
	)
}
