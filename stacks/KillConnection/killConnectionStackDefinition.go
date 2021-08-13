package KillConnection

import (
	"context"
	"github.com/bhbosman/gocomms/internal"
	"github.com/bhbosman/goerrors"
	"github.com/bhbosman/goprotoextra"
	"github.com/google/uuid"
	"github.com/reactivex/rxgo/v2"
)

const StackName = "KillConnection"

func StackDefinition(
	connectionType internal.ConnectionType,
	cancelContext context.Context,
	stackCancelFunc internal.CancelFunc,
	connectionManager rxgo.IPublishToConnectionManager,
	opts ...rxgo.Option) (*internal.StackDefinition, error) {
	if stackCancelFunc == nil {
		return nil, goerrors.InvalidParam
	}
	id := uuid.New()
	return &internal.StackDefinition{
		IId:  id,
		Name: StackName,
		Inbound: internal.NewBoundResultImpl(func(inOutBoundParams internal.InOutBoundParams) (internal.IStackBoundDefinition, error) {
			return &internal.StackBoundDefinition{
					PipeDefinition: func(stackData, pipeData interface{}, pipeParams internal.PipeDefinitionParams) (uuid.UUID, rxgo.Observable, error) {
						return id,
							pipeParams.Obs.(rxgo.InOutBoundObservable).MapInOutBound(
								inOutBoundParams.Index,
								pipeParams.ConnectionId,
								StackName,
								rxgo.StreamDirectionInbound,
								pipeParams.ConnectionManager,
								func(ctx context.Context, i goprotoextra.ReadWriterSize) (goprotoextra.ReadWriterSize, error) {
									return i, nil
								},
								opts...), nil
					},
				},
				nil
		}),
		Outbound: internal.NewBoundResultImpl(func(inOutBoundParams internal.InOutBoundParams) (internal.IStackBoundDefinition, error) {
			var outBoundChannel *internal.ChannelManager
			return &internal.StackBoundDefinition{
					PipeDefinition: func(stackData, pipeData interface{}, pipeParams internal.PipeDefinitionParams) (uuid.UUID, rxgo.Observable, error) {
						outBoundChannel = internal.NewChannelManager("", "")
						_ = pipeParams.Obs.(rxgo.InOutBoundObservable).DoOnNextInOutBound(
							inOutBoundParams.Index,
							pipeParams.ConnectionId,
							StackName,
							rxgo.StreamDirectionOutbound,
							connectionManager,
							func(ctx context.Context, rws goprotoextra.ReadWriterSize) {
								outBoundChannel.Send(ctx, rws)
							}, opts...)
						result := rxgo.FromChannel(outBoundChannel.Items, opts...)
						return id, result, nil
					},
				},
				nil
		}),
	}, nil
}
