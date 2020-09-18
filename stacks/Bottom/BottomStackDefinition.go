package Bottom

import (
	"context"
	rxgo "github.com/ReactiveX/RxGo"
	"github.com/bhbosman/gocommon/stacks/defs"
	"github.com/bhbosman/goerrors"
	"github.com/bhbosman/goprotoextra"
	"net"
	"net/url"
)

func StackDefinition() (*defs.StackDefinition, error) {

	return &defs.StackDefinition{
		Name: goerrors.BottomStackName,
		Inbound: func(index int, ctx context.Context) defs.BoundDefinition {
			return defs.BoundDefinition{
				PipeDefinition: func(params defs.PipeDefinitionParams) (rxgo.Observable, error) {
					return params.Obs.(rxgo.InOutBoundObservable).MapInOutBound(
						index,
						params.ConnectionId,
						goerrors.BottomStackName,
						rxgo.StreamDirectionInbound,
						params.ConnectionManager,
						func(ctx context.Context, rws goprotoextra.ReadWriterSize) (goprotoextra.ReadWriterSize, error) {
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
						goerrors.BottomStackName,
						rxgo.StreamDirectionOutbound,
						params.ConnectionManager,
						func(ctx context.Context, rws goprotoextra.ReadWriterSize) (goprotoextra.ReadWriterSize, error) {
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
