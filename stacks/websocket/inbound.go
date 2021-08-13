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
	connectionType internal.ConnectionType,
	stackCancelFunc internal.CancelFunc,
	connectionManager rxgo.IPublishToConnectionManager, opts ...rxgo.Option) internal.BoundResult {
	return func(inOutBoundParams internal.InOutBoundParams) (internal.IStackBoundDefinition, error) {
		return &internal.StackBoundDefinition{
			PipeDefinition: func(stackData, pipeData interface{}, pipeParams internal.PipeDefinitionParams) (uuid.UUID, rxgo.Observable, error) {
				if stackCancelFunc == nil {
					return uuid.Nil, nil, goerrors.InvalidParam
				}
				sd, ok := stackData.(*StackData)
				if !ok {
					return uuid.Nil, nil, WrongStackDataError(connectionType, stackData)
				}
				_ = pipeParams.Obs.(rxgo.InOutBoundObservable).DoOnNextInOutBound(
					inOutBoundParams.Index,
					pipeParams.ConnectionId,
					StackName,
					rxgo.StreamDirectionInbound,
					connectionManager,
					func(ctx context.Context, rws goprotoextra.ReadWriterSize) {
						_, err := io.Copy(sd.pipeWriteClose, rws)
						if err != nil {
							return
						}
					},
					opts...)
				return id, rxgo.FromChannel(sd.nextInBoundChannelManager.Items), nil
			},
			//PipeState: &internal.PipeState{
			//	Destroy: func(stackData interface{}) interface{} {
			//		if closer, ok := stackData.(io.Closer); ok {
			//			return closer.Close()
			//		}
			//		return nil
			//	},
			//	Start: func(stackData, pipeData interface{}, ctx context.Context) error {
			//		return ctx.Err()
			//	},
			//	ID: id,
			//	Create: func(stackData interface{}, ctx context.Context) (interface{}, error) {
			//		return internal.NewNoCloser(), nil
			//	},
			//	End: func(stackData, pipeData interface{}) error {
			//		var err error = nil
			//		if closer, ok := pipeData.(io.Closer); ok {
			//			err = multierr.Append(err, closer.Close())
			//		}
			//		return err
			//	},
			//},
		}, nil
	}
}
