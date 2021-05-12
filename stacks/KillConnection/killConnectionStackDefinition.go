package KillConnection

import (
	"context"
	"github.com/bhbosman/gocomms/internal"
	"github.com/bhbosman/goerrors"
	"github.com/bhbosman/gomessageblock"
	"github.com/bhbosman/goprotoextra"
	"github.com/google/uuid"
	"github.com/reactivex/rxgo/v2"
	"time"
)

func StackDefinition(
	connectionType internal.ConnectionType,
	cancelContext context.Context,
	stackCancelFunc internal.CancelFunc,
	connectionManager rxgo.IPublishToConnectionManager,
	opts ...rxgo.Option) (*internal.StackDefinition, error) {
	if stackCancelFunc == nil {
		return nil, goerrors.InvalidParam
	}
	const stackName = "KillConnection"
	id := uuid.New()
	return &internal.StackDefinition{
		Id:   id,
		Name: stackName,
		Inbound: internal.NewBoundResultImpl(
			func(stackData, pipeData interface{}, inOutBoundParams internal.InOutBoundParams) internal.IStackBoundDefinition {
				return &internal.StackBoundDefinition{
					PipeDefinition: func(
						pipeParams internal.PipeDefinitionParams) (uuid.UUID, rxgo.Observable, error) {
						if stackCancelFunc == nil {
							return uuid.Nil, nil, goerrors.InvalidParam
						}
						return id,
							pipeParams.Obs.(rxgo.InOutBoundObservable).MapInOutBound(
								inOutBoundParams.Index,
								pipeParams.ConnectionId,
								stackName,
								rxgo.StreamDirectionInbound,
								pipeParams.ConnectionManager,
								func(ctx context.Context, i goprotoextra.ReadWriterSize) (goprotoextra.ReadWriterSize, error) {
									return i, nil
								},
								opts...), nil
					},
				}
			}),
		Outbound: internal.NewBoundResultImpl(
			func(stackData, pipeData interface{}, inOutBoundParams internal.InOutBoundParams) internal.IStackBoundDefinition {
				var outBoundChannel chan rxgo.Item
				return &internal.StackBoundDefinition{
					PipeDefinition: func(pipeParams internal.PipeDefinitionParams) (uuid.UUID, rxgo.Observable, error) {
						if stackCancelFunc == nil {
							return uuid.Nil, nil, goerrors.InvalidParam
						}
						outBoundChannel = make(chan rxgo.Item)
						_ = pipeParams.Obs.(rxgo.InOutBoundObservable).DoOnNextInOutBound(
							inOutBoundParams.Index,
							pipeParams.ConnectionId,
							stackName,
							rxgo.StreamDirectionOutbound,
							connectionManager,
							func(ctx context.Context, rws goprotoextra.ReadWriterSize) {
								outBoundChannel <- rxgo.Of(rws)
							}, opts...)
						result := rxgo.FromChannel(outBoundChannel, opts...)
						return id, result, nil
					},
					PipeState: &internal.PipeState{
						Start: func(ctx context.Context) error {
							if cancelContext.Err() == nil {
								go func() {
									outBoundChannel <- rxgo.Of(gomessageblock.NewReaderWriterBlock([]byte("ERR:No Transport layer selected. Closing down connection\n")))
									time.Sleep(time.Millisecond * 10)
									stackCancelFunc("Kill Connection", false, goerrors.InvalidParam)
									return
								}()
							}
							return nil
						},
						End: func() error {
							close(outBoundChannel)
							return nil
						},
					},
				}
			}),
	}, nil
}
