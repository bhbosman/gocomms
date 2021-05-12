package internal

import (
	"context"
	"github.com/bhbosman/gocomms/internal"
	"github.com/bhbosman/goprotoextra"
	"github.com/google/uuid"
	"github.com/reactivex/rxgo/v2"
)

func Outbound(id uuid.UUID, opts ...rxgo.Option) internal.BoundResult {
	return func(stackData, pipeData interface{}, inOutBoundParams internal.InOutBoundParams) internal.IStackBoundDefinition {
		return &internal.StackBoundDefinition{
			PipeDefinition: func(pipeParams internal.PipeDefinitionParams) (uuid.UUID, rxgo.Observable, error) {
				return id,
					pipeParams.Obs.(rxgo.InOutBoundObservable).MapInOutBound(
						inOutBoundParams.Index,
						pipeParams.ConnectionId,
						StackName,
						rxgo.StreamDirectionOutbound,
						pipeParams.ConnectionManager,
						func(ctx context.Context, rws goprotoextra.ReadWriterSize) (goprotoextra.ReadWriterSize, error) {
							return rws, nil
						},
						opts...),
					nil
			},
		}
	}
}
