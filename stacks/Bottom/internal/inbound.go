package internal

import (
	"context"
	"github.com/bhbosman/gocomms/internal"
	"github.com/bhbosman/goprotoextra"
	"github.com/google/uuid"
	"github.com/reactivex/rxgo/v2"
)

func Inbound(id uuid.UUID, opts ...rxgo.Option) internal.BoundResult {
	return func(inOutBoundParams internal.InOutBoundParams) (internal.IStackBoundDefinition, error) {
		return &internal.StackBoundDefinition{
				PipeDefinition: func(stackData, pipeData interface{}, pipeParams internal.PipeDefinitionParams) (uuid.UUID, rxgo.Observable, error) {
					return id, pipeParams.Obs.(rxgo.InOutBoundObservable).MapInOutBound(
							inOutBoundParams.Index,
							pipeParams.ConnectionId,
							StackName,
							rxgo.StreamDirectionInbound,
							pipeParams.ConnectionManager,
							func(ctx context.Context, rws goprotoextra.ReadWriterSize) (goprotoextra.ReadWriterSize, error) {
								return rws, nil
							},
							opts...),
						nil
				},
			},
			nil
	}
}
