package websocket

import (
	"context"
	"github.com/bhbosman/gocommon/stream"
	"github.com/bhbosman/gocomms/internal"
	"github.com/bhbosman/gocomms/stacks/websocket/wsmsg"
	"github.com/bhbosman/goerrors"
	"github.com/bhbosman/gomessageblock"
	"github.com/bhbosman/goprotoextra"
	"github.com/gobwas/ws"
	"github.com/gobwas/ws/wsutil"
	"github.com/google/uuid"
	"github.com/reactivex/rxgo/v2"
	"time"
)

func outbound(
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
					inOutBoundParams.Index-1,
					pipeParams.ConnectionId,
					StackName+"FromUser",
					rxgo.StreamDirectionOutbound,
					connectionManager,
					func(ctx context.Context, size goprotoextra.ReadWriterSize) {
						data.tempStep.Send(ctx, size)
					},
					opts...)

				tempObs := rxgo.FromChannel(data.tempStep.Items).(rxgo.InOutBoundObservable)
				_ = tempObs.DoOnNextInOutBound(
					inOutBoundParams.Index,
					pipeParams.ConnectionId,
					StackName,
					rxgo.StreamDirectionOutbound,
					connectionManager,
					func(ctx context.Context, size goprotoextra.ReadWriterSize) {
						switch v := size.(type) {
						case *gomessageblock.ReaderWriter:
							messageWrapper, err := stream.UnMarshal(v, nil, nil, nil, nil)
							if err != nil {
								return
							}
							switch vv := messageWrapper.(type) {
							case *wsmsg.WebSocketMessageWrapper:
								switch vv.Data.OpCode {
								case wsmsg.WebSocketMessage_OpText:
									err := wsutil.WriteClientMessage(data.upgradedConnection, ws.OpText, vv.Data.Message)
									if err != nil {
										stackCancelFunc("creating text payload", false, err)
										return
									}
								case wsmsg.WebSocketMessage_OpPing:
									binary, _ := time.Now().MarshalBinary()
									err := wsutil.WriteClientMessage(data.upgradedConnection, ws.OpPing, binary)
									if err != nil {
										stackCancelFunc("creating ping payload", false, err)
										return
									}
								}
							}
						}
					},
					opts...)
				return id, rxgo.FromChannel(data.nextOutboundChannelManager.Items), nil
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
