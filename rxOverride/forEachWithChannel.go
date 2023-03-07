package rxOverride

import (
	"context"
	"fmt"
	"github.com/bhbosman/goCommsDefinitions"
	"github.com/bhbosman/gocommon/GoFunctionCounter"
	"github.com/bhbosman/gocommon/model"
	"github.com/bhbosman/gocommon/rxHelper"
	"github.com/reactivex/rxgo/v2"
)

func ForEachWithChannel(
	description string,
	direction model.StreamDirection,
	o rxgo.Observable,
	ctx context.Context,
	goFunctionCounter GoFunctionCounter.IService,
	channel <-chan interface{},
	handler goCommsDefinitions.IRxNextHandler,
	opts ...rxgo.Option,
) rxgo.Disposed {
	dispose := make(chan struct{})
	localHandler := func(ctx context.Context, src *rxgo.ItemChannel) {
		defer close(dispose)

		localNextFunc := rxHelper.HandleMessage(
			handler.OnSendData,
			handler.OnTrySendData,
			handler.IsActive,
			func() int {
				return src.Len()
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
				localNextFunc(i)
				break
			case i, ok := <-src.Ch:
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
