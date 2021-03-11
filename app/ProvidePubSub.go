package app

import (
	"context"
	"github.com/cskr/pubsub"
	"go.uber.org/fx"
)

func ProvidePubSub(name string, pubSub *pubsub.PubSub) fx.Option {
	return fx.Provide(
		fx.Annotated{
			Name: name,
			Target: func(lifecycle fx.Lifecycle) *pubsub.PubSub {
				result := pubSub
				lifecycle.Append(fx.Hook{
					OnStart: nil,
					OnStop: func(ctx context.Context) error {
						result.Shutdown()
						return nil
					},
				})
				return result
			},
		},
	)
}
