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

func ForEach2(
	description string,
	direction model.StreamDirection,
	o rxgo.Observable,
	ctx context.Context,
	goFunctionCounter GoFunctionCounter.IService,
	handler goCommsDefinitions.IRxNextHandler,
	opts ...rxgo.Option,
) rxgo.Disposed {
	return ForEach1(
		description,
		direction, o, ctx, goFunctionCounter,
		handler.OnSendData,
		handler.OnTrySendData,
		handler.OnError,
		handler.OnComplete,
		handler.IsActive,
		opts...)
}

func ForEach1(
	description string,
	direction model.StreamDirection,
	o rxgo.Observable,
	ctx context.Context,
	goFunctionCounter GoFunctionCounter.IService,
	nextFunc rxgo.NextFunc,
	tryNextFunc goCommsDefinitions.TryNextFunc,
	errFunc rxgo.ErrFunc,
	completedFunc rxgo.CompletedFunc,
	isActive func() bool,
	opts ...rxgo.Option,
) rxgo.Disposed {
	dispose := make(chan struct{})
	handler := func(ctx context.Context, src <-chan rxgo.Item) {
		defer close(dispose)

		localNextFunc := rxHelper.HandleMessage(
			nextFunc,
			tryNextFunc,
			isActive,
			func() int {
				return len(src)
			},
		)
		for {
			select {
			case <-ctx.Done():
				completedFunc()
				return
			case i, ok := <-src:
				if !ok {
					completedFunc()
					return
				}
				if i.Error() {
					errFunc(i.E)
					break
				}
				localNextFunc(i.V)
			}
		}
	}
	_ = goFunctionCounter.GoRun(
		fmt.Sprintf("foreach.%v.%v", description, direction.String()),
		func() {
			handler(ctx, o.Observe(opts...))
		},
	)
	return dispose
}
