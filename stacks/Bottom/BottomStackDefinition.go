package Bottom

import (
	"context"
	"github.com/bhbosman/gocomms/internal"
	"github.com/bhbosman/goerrors"
	"github.com/bhbosman/goprotoextra"
	"github.com/reactivex/rxgo/v2"
	"net"
	"net/url"
)

func StackDefinition(opts ...rxgo.Option) (*internal.StackDefinition, error) {
	return &internal.StackDefinition{
		Name: goerrors.BottomStackName,
		Inbound: func(index int, ctx context.Context) internal.BoundDefinition {
			return internal.BoundDefinition{
				PipeDefinition: func(params internal.PipeDefinitionParams) (rxgo.Observable, error) {
					return params.Obs.(rxgo.InOutBoundObservable).MapInOutBound(
						index,
						params.ConnectionId,
						goerrors.BottomStackName,
						rxgo.StreamDirectionInbound,
						params.ConnectionManager,
						func(ctx context.Context, rws goprotoextra.ReadWriterSize) (goprotoextra.ReadWriterSize, error) {
							return rws, nil
						},
						opts...), nil
				},
			}
		},
		Outbound: func(index int, ctx context.Context) internal.BoundDefinition {
			return internal.BoundDefinition{
				PipeDefinition: func(params internal.PipeDefinitionParams) (rxgo.Observable, error) {
					return params.Obs.(rxgo.InOutBoundObservable).MapInOutBound(
						index,
						params.ConnectionId,
						goerrors.BottomStackName,
						rxgo.StreamDirectionOutbound,
						params.ConnectionManager,
						func(ctx context.Context, rws goprotoextra.ReadWriterSize) (goprotoextra.ReadWriterSize, error) {
							return rws, nil
						},
						opts...), nil
				},
			}
		},
		StackState: internal.StackState{
			Start: func(conn net.Conn, url *url.URL, ctx context.Context, cancelFunc internal.CancelFunc) (net.Conn, error) {
				return conn, ctx.Err()
			},
		},
	}, nil
}
