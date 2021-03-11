package Top

import (
	"context"
	"github.com/bhbosman/gocomms/internal"
	"github.com/bhbosman/goerrors"
	"github.com/bhbosman/goprotoextra"
	"github.com/reactivex/rxgo/v2"
)

func StackDefinition(
	connectionType internal.ConnectionType,
	opts ...rxgo.Option) (*internal.StackDefinition, error) {
	return &internal.StackDefinition{
		Name: goerrors.TopStackName,
		Inbound: func(inOutBoundParams internal.InOutBoundParams) internal.BoundDefinition {
			return internal.BoundDefinition{
				PipeDefinition: func(pipeParams internal.PipeDefinitionParams) (rxgo.Observable, error) {
					return pipeParams.Obs.(rxgo.InOutBoundObservable).MapInOutBound(
						inOutBoundParams.Index,
						pipeParams.ConnectionId,
						goerrors.TopStackName,
						rxgo.StreamDirectionInbound,
						pipeParams.ConnectionManager,
						func(ctx context.Context, rws goprotoextra.ReadWriterSize) (goprotoextra.ReadWriterSize, error) {
							return rws, ctx.Err()
						}), nil
				},
			}
		},
		Outbound: func(inOutBoundParams internal.InOutBoundParams) internal.BoundDefinition {
			return internal.BoundDefinition{
				PipeDefinition: func(pipeParams internal.PipeDefinitionParams) (rxgo.Observable, error) {
					return pipeParams.Obs.(rxgo.InOutBoundObservable).MapInOutBound(
						inOutBoundParams.Index,
						pipeParams.ConnectionId,
						goerrors.TopStackName,
						rxgo.StreamDirectionOutbound,
						pipeParams.ConnectionManager,
						func(ctx context.Context, rws goprotoextra.ReadWriterSize) (goprotoextra.ReadWriterSize, error) {
							return rws, ctx.Err()
						}), nil
				},
			}
		},
	}, nil
}
