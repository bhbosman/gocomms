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
		Id:   id,
		Name: "CrMessageBreaker",
		Inbound: internal.NewBoundResultImpl(
			func(stackData, pipeData interface{}, inOutBoundParams internal2.InOutBoundParams) internal2.IStackBoundDefinition {
				return &internal.StackBoundDefinition{
					PipeDefinition: func(pipeParams internal.PipeDefinitionParams) (uuid.UUID, rxgo.Observable, error) {

						channelManager := internal2.NewChannelManager(make(chan rxgo.Item), "Inbound MessageBreaker", connectionId)
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
				}
			}),
		Outbound: internal.NewBoundResultImpl(
			func(stackData, pipeData interface{}, inOutBoundParams internal2.InOutBoundParams) internal2.IStackBoundDefinition {
				return &internal.StackBoundDefinition{
					PipeDefinition: func(pipeParams internal.PipeDefinitionParams) (uuid.UUID, rxgo.Observable, error) {
						return id,
							pipeParams.Obs.(rxgo.InOutBoundObservable).MapInOutBound(
								inOutBoundParams.Index,
								pipeParams.ConnectionId,
								"CrMessageBreaker",
								rxgo.StreamDirectionInbound,
								pipeParams.ConnectionManager,
								func(ctx context.Context, rws goprotoextra.ReadWriterSize) (goprotoextra.ReadWriterSize, error) {
									return rws, nil
								},
								opts...),
							nil
					},
				}
			}),
		StackState: internal2.StackState{
			Start: nil,
			End:   nil,
		},
	}, nil
}
