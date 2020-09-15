package Top

import (
	"context"
	rxgo "github.com/ReactiveX/RxGo"
	"github.com/bhbosman/gocommon/constants"
	"github.com/bhbosman/gocommon/stacks/defs"
	"github.com/bhbosman/gocommon/stream"
)

func StackDefinition() (*defs.StackDefinition, error) {
	return &defs.StackDefinition{
		Name: constants.TopStackName,
		Inbound: func(index int, ctx context.Context) defs.BoundDefinition {
			return defs.BoundDefinition{
				PipeDefinition: func(params defs.PipeDefinitionParams) (rxgo.Observable, error) {
					return params.Obs.(rxgo.InOutBoundObservable).MapInOutBound(
						index,
						params.ConnectionId,
						constants.TopStackName,
						rxgo.StreamDirectionInbound,
						params.ConnectionManager,
						func(ctx context.Context, rws stream.ReadWriterSize) (stream.ReadWriterSize, error) {
							return rws, ctx.Err()
						}), nil
				},
			}
		},
		Outbound: func(index int, ctx context.Context) defs.BoundDefinition {
			return defs.BoundDefinition{
				PipeDefinition: func(params defs.PipeDefinitionParams) (rxgo.Observable, error) {
					return params.Obs.(rxgo.InOutBoundObservable).MapInOutBound(
						index,
						params.ConnectionId,
						constants.TopStackName,
						rxgo.StreamDirectionOutbound,
						params.ConnectionManager,
						func(ctx context.Context, rws stream.ReadWriterSize) (stream.ReadWriterSize, error) {
							return rws, ctx.Err()
						}), nil
				},
			}
		},
	}, nil
}
