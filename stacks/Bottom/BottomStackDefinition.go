package Bottom

import (
	"context"
	rxgo "github.com/ReactiveX/RxGo"
	"github.com/bhbosman/gocommon/constants"
	"github.com/bhbosman/gocommon/stacks/defs"
	"github.com/bhbosman/gocommon/stream"
	"net"
	"net/url"
)

func StackDefinition() (*defs.StackDefinition, error) {

	return &defs.StackDefinition{
		Name: constants.BottomStackName,
		Inbound: func(index int, ctx context.Context) defs.BoundDefinition {
			return defs.BoundDefinition{
				PipeDefinition: func(params defs.PipeDefinitionParams) (rxgo.Observable, error) {
					return params.Obs.(rxgo.InOutBoundObservable).MapInOutBound(
						index,
						params.ConnectionId,
						constants.BottomStackName,
						rxgo.StreamDirectionInbound,
						params.ConnectionManager,
						func(ctx context.Context, rws stream.ReadWriterSize) (stream.ReadWriterSize, error) {
							return rws, nil
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
						constants.BottomStackName,
						rxgo.StreamDirectionOutbound,
						params.ConnectionManager,
						func(ctx context.Context, rws stream.ReadWriterSize) (stream.ReadWriterSize, error) {
							return rws, nil
						}), nil
				},
			}
		},
		StackState: defs.StackState{
			Start: func(conn net.Conn, url *url.URL, ctx context.Context, cancelFunc defs.CancelFunc) (net.Conn, error) {
				return conn, ctx.Err()
			},
		},
	}, nil
}
