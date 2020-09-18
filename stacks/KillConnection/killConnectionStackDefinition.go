package KillConnection

import (
	"context"
	"github.com/bhbosman/gocommon/multiBlock"
	"github.com/bhbosman/gocommon/stacks/defs"
	"github.com/bhbosman/goerrors"
	"github.com/bhbosman/goprotoextra"
	"github.com/reactivex/rxgo/v2"
	"time"
)

func StackDefinition(
	cancelContext context.Context,
	stackCancelFunc defs.CancelFunc,
	connectionManager rxgo.IPublishToConnectionManager,
	opts ...rxgo.Option) (*defs.StackDefinition, error) {
	if stackCancelFunc == nil {
		return nil, goerrors.InvalidParam
	}
	const stackName = "KillConnection"
	return &defs.StackDefinition{
		Name: stackName,
		Inbound: func(index int, ctx context.Context) defs.BoundDefinition {
			return defs.BoundDefinition{
				PipeDefinition: func(
					params defs.PipeDefinitionParams) (rxgo.Observable, error) {
					if stackCancelFunc == nil {
						return nil, goerrors.InvalidParam
					}
					return params.Obs.(rxgo.InOutBoundObservable).MapInOutBound(
						index,
						params.ConnectionId,
						stackName,
						rxgo.StreamDirectionInbound,
						params.ConnectionManager,
						func(ctx context.Context, i goprotoextra.ReadWriterSize) (goprotoextra.ReadWriterSize, error) {
							return i, nil
						},
						opts...), nil
				},
			}
		},
		Outbound: func(index int, ctx context.Context) defs.BoundDefinition {
			var outBoundChannel chan rxgo.Item
			return defs.BoundDefinition{
				PipeDefinition: func(params defs.PipeDefinitionParams) (rxgo.Observable, error) {
					if stackCancelFunc == nil {
						return nil, goerrors.InvalidParam
					}
					outBoundChannel = make(chan rxgo.Item)
					_ = params.Obs.(rxgo.InOutBoundObservable).DoOnNextInOutBound(
						index,
						params.ConnectionId,
						stackName,
						rxgo.StreamDirectionOutbound,
						connectionManager,
						func(ctx context.Context, rws goprotoextra.ReadWriterSize) {
							outBoundChannel <- rxgo.Of(rws)
						})
					result := rxgo.FromChannel(outBoundChannel, opts...)
					return result, nil
				},
				PipeState: defs.PipeState{
					Start: func(ctx context.Context) error {
						if cancelContext.Err() == nil {
							go func() {
								outBoundChannel <- rxgo.Of(multiBlock.NewReaderWriterBlock([]byte("ERR:No Transport layer selected. Closing down connection\n")))
								time.Sleep(time.Millisecond * 10)
								stackCancelFunc("Kill Connection", false, goerrors.InvalidParam)
								return
							}()
						}
						return nil
					},
					End: func() error {
						close(outBoundChannel)
						return nil
					},
				},
			}
		},
	}, nil
}
