package websocket

import (
	"context"
	"github.com/bhbosman/gocomms/internal"
	"github.com/bhbosman/goerrors"
	"github.com/bhbosman/goprotoextra"
	"github.com/google/uuid"
	"github.com/reactivex/rxgo/v2"
	"io"
)

func inbound(
	id uuid.UUID,
	data *Data,
	stackCancelFunc internal.CancelFunc,
	connectionManager rxgo.IPublishToConnectionManager, opts ...rxgo.Option) internal.BoundResult {
	return func(stackData, pipeData interface{}, inOutBoundParams internal.InOutBoundParams) internal.IStackBoundDefinition {
		return &internal.StackBoundDefinition{
			PipeDefinition: func(pipeParams internal.PipeDefinitionParams) (uuid.UUID, rxgo.Observable, error) {
				if stackCancelFunc == nil {
					return uuid.Nil, nil, goerrors.InvalidParam
				}
				_ = pipeParams.Obs.(rxgo.InOutBoundObservable).DoOnNextInOutBound(
					inOutBoundParams.Index,
					pipeParams.ConnectionId,
					StackName,
					rxgo.StreamDirectionInbound,
					connectionManager,
					func(ctx context.Context, rws goprotoextra.ReadWriterSize) {
						_, err := io.Copy(data.pipeWriteClose, rws)
						if err != nil {
							return
						}
					},
					opts...)
				return id, rxgo.FromChannel(data.nextInBoundChannelManager.Items), nil
			},
			PipeState: &internal.PipeState{
				Start: func(ctx context.Context) error {
					return ctx.Err()
				},
				End: func() error {
					return nil
				},
			},
		}
	}
}
