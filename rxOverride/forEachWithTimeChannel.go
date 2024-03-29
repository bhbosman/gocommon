package rxOverride

import (
	"context"
	"fmt"
	"github.com/bhbosman/goCommsDefinitions"
	"github.com/bhbosman/gocommon"
	"github.com/bhbosman/gocommon/GoFunctionCounter"
	"github.com/bhbosman/gocommon/model"
	"github.com/bhbosman/gocommon/rxHelper"
	"github.com/reactivex/rxgo/v2"
	"time"
)

func ForEachWithTimeChannel(
	description string,
	direction model.StreamDirection,
	o gocommon.IObservable,
	ctx context.Context,
	goFunctionCounter GoFunctionCounter.IService,
	handler goCommsDefinitions.IRxNextHandler,
	channel <-chan time.Time,
	channelHandler func(time2 time.Time) interface{},
	opts ...rxgo.Option,
) rxgo.Disposed {
	dispose := make(chan struct{})
	localHandler := func(ctx context.Context, src <-chan rxgo.Item) {
		defer close(dispose)

		localNextFunc := rxHelper.HandleMessage(
			handler.OnSendData,
			handler.OnTrySendData,
			handler.IsActive,
			func() int {
				return len(src)
			},
			func() int {
				return len(channel)
			},
		)
		for {
			select {
			case <-ctx.Done():
				handler.OnComplete()
				return
			case i, ok := <-channel:
				if !ok {
					handler.OnComplete()
					return
				}
				var unk interface{} = i
				if channelHandler != nil {
					unk = channelHandler(i)
				}
				localNextFunc(unk)
				break
			case i, ok := <-src:
				if !ok {
					handler.OnComplete()
					return
				}
				if i.Error() {
					handler.OnError(i.E)
					break
				}
				localNextFunc(i.V)
			}
		}
	}
	_ = goFunctionCounter.GoRun(
		fmt.Sprintf("foreach.%v.%v", description, direction.String()),
		func() {
			localHandler(ctx, o.Observe(opts...))
		},
	)
	return dispose
}
