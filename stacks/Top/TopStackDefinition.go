package Top

import (
	"context"
	"github.com/bhbosman/gocomms/internal"
	"github.com/bhbosman/goerrors"
	"github.com/bhbosman/goprotoextra"
	"github.com/reactivex/rxgo/v2"
)

func StackDefinition() (*internal.StackDefinition, error) {
	return &internal.StackDefinition{
		Name: goerrors.TopStackName,
		Inbound: func(index int, ctx context.Context) internal.BoundDefinition {
			return internal.BoundDefinition{
				PipeDefinition: func(params internal.PipeDefinitionParams) (rxgo.Observable, error) {
					return params.Obs.(rxgo.InOutBoundObservable).MapInOutBound(
						index,
						params.ConnectionId,
						goerrors.TopStackName,
						rxgo.StreamDirectionInbound,
						params.ConnectionManager,
						func(ctx context.Context, rws goprotoextra.ReadWriterSize) (goprotoextra.ReadWriterSize, error) {
							return rws, ctx.Err()
						}), nil
				},
			}
		},
		Outbound: func(index int, ctx context.Context) internal.BoundDefinition {
			return internal.BoundDefinition{
				PipeDefinition: func(params internal.PipeDefinitionParams) (rxgo.Observable, error) {
					return params.Obs.(rxgo.InOutBoundObservable).MapInOutBound(
						index,
						params.ConnectionId,
						goerrors.TopStackName,
						rxgo.StreamDirectionOutbound,
						params.ConnectionManager,
						func(ctx context.Context, rws goprotoextra.ReadWriterSize) (goprotoextra.ReadWriterSize, error) {
							return rws, ctx.Err()
						}), nil
				},
			}
		},
	}, nil
}
