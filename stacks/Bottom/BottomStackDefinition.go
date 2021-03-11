package Bottom

import (
	"context"
	"github.com/bhbosman/gocomms/internal"
	"github.com/bhbosman/goerrors"
	"github.com/bhbosman/goprotoextra"
	"github.com/reactivex/rxgo/v2"
	"net"
)

func StackDefinition(
	connectionType internal.ConnectionType,
	opts ...rxgo.Option) (*internal.StackDefinition, error) {
	return &internal.StackDefinition{
		Name: goerrors.BottomStackName,
		Inbound: func(inOutBoundParams internal.InOutBoundParams) internal.BoundDefinition {
			return internal.BoundDefinition{
				PipeDefinition: func(pipeParams internal.PipeDefinitionParams) (rxgo.Observable, error) {
					return pipeParams.Obs.(rxgo.InOutBoundObservable).MapInOutBound(
						inOutBoundParams.Index,
						pipeParams.ConnectionId,
						goerrors.BottomStackName,
						rxgo.StreamDirectionInbound,
						pipeParams.ConnectionManager,
						func(ctx context.Context, rws goprotoextra.ReadWriterSize) (goprotoextra.ReadWriterSize, error) {
							return rws, nil
						},
						opts...), nil
				},
			}
		},
		Outbound: func(inOutBoundParams internal.InOutBoundParams) internal.BoundDefinition {
			return internal.BoundDefinition{
				PipeDefinition: func(pipeParams internal.PipeDefinitionParams) (rxgo.Observable, error) {
					return pipeParams.Obs.(rxgo.InOutBoundObservable).MapInOutBound(
						inOutBoundParams.Index,
						pipeParams.ConnectionId,
						goerrors.BottomStackName,
						rxgo.StreamDirectionOutbound,
						pipeParams.ConnectionManager,
						func(ctx context.Context, rws goprotoextra.ReadWriterSize) (goprotoextra.ReadWriterSize, error) {
							return rws, nil
						},
						opts...), nil
				},
			}
		},
		StackState: internal.StackState{
			Start: func(startParams internal.StackStartStateParams) (net.Conn, error) {
				return startParams.Conn, startParams.Ctx.Err()
			},
		},
	}, nil
}
