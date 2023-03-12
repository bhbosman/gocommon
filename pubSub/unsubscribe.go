package pubSub

import (
	"github.com/bhbosman/goCommsDefinitions"
	"github.com/cskr/pubsub"
)

type IPubSub interface {
	Pub(msg interface{}, topics ...string)
	PubWithContext(msg interface{}, topics ...string) bool
}

func Unsubscribe(
	name string,
	pubSub *pubsub.PubSub,
	unk interface {
		GoRun(s string, cb func()) error
	},
	subscribeChannel goCommsDefinitions.IPubSubBag,
) error {
	_ = unk.GoRun(
		name,
		func() {
			pubSub.Unsub(subscribeChannel)
		},
	)
	_ = unk.GoRun(
		name,
		func() {
			if subscribeChannel != nil {
				subscribeChannel.Flush()
			}
		},
	)
	return nil
}
