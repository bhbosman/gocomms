package crMessageBreaker

import (
	"context"
	"github.com/bhbosman/gocomms/internal"
	internal2 "github.com/bhbosman/gocomms/internal"
	"github.com/bhbosman/goprotoextra"
	"github.com/google/uuid"
	"github.com/reactivex/rxgo/v2"
)

func StackDefinition(
	connectionType internal2.ConnectionType,
	connectionId string,
	stackCancelFunc internal2.CancelFunc,
	connectionManager rxgo.IPublishToConnectionManager,
	opts ...rxgo.Option) (*internal2.StackDefinition, error) {
	id := uuid.New()
	return &internal2.StackDefinition{
		IId:  id,
		Name: StackName,
		Inbound: internal.NewBoundResultImpl(func(inOutBoundParams internal2.InOutBoundParams) (internal2.IStackBoundDefinition, error) {
			return &internal.StackBoundDefinition{
				PipeDefinition: func(stackData, pipeData interface{}, pipeParams internal.PipeDefinitionParams) (uuid.UUID, rxgo.Observable, error) {
					channelManager := internal2.NewChannelManager("Inbound MessageBreaker", connectionId)
					disposable := pipeParams.Obs.(rxgo.InOutBoundObservable).DoOnNextInOutBound(
						inOutBoundParams.Index,
						pipeParams.ConnectionId,
						"CrMessageBreaker",
						rxgo.StreamDirectionInbound,
						pipeParams.ConnectionManager,
						func(ctx context.Context, size goprotoextra.ReadWriterSize) {

						},
						opts...)
					go func() {
						<-disposable
						_ = channelManager.Close()
					}()

					return id, rxgo.FromChannel(channelManager.Items, opts...), nil
				},
			}, nil
		}),
		Outbound: internal.NewBoundResultImpl(func(inOutBoundParams internal2.InOutBoundParams) (internal2.IStackBoundDefinition, error) {
			return &internal.StackBoundDefinition{
				PipeDefinition: func(stackData, pipeData interface{}, pipeParams internal.PipeDefinitionParams) (uuid.UUID, rxgo.Observable, error) {
					return id,
						pipeParams.Obs.(rxgo.InOutBoundObservable).MapInOutBound(
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
			}, nil
		}),
	}, nil
}
