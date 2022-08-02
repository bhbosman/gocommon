package PubSub

import (
	"github.com/cskr/pubsub"
	"go.uber.org/fx"
)

func ProvidePubSubInstance(name string, sub *pubsub.PubSub) fx.Option {
	return fx.Provide(
		fx.Annotated{
			Name: name,
			Target: func() (*pubsub.PubSub, error) {
				return sub, nil
			},
		},
	)

}
