package websocket

import (
	"context"
	"fmt"
	"github.com/bhbosman/gocommon/stream"
	"github.com/bhbosman/gocomms/internal"
	"github.com/bhbosman/gocomms/stacks/internal/connectionWrapper"
	"github.com/bhbosman/gocomms/stacks/websocket/wsmsg"
	"github.com/bhbosman/goerrors"
	"github.com/bhbosman/gomessageblock"
	"github.com/bhbosman/goprotoextra"
	"github.com/gobwas/ws"
	"github.com/gobwas/ws/wsutil"
	"github.com/reactivex/rxgo/v2"
	"go.uber.org/multierr"
	"io"
	"net"
	"net/url"
	"time"
)

func StackDefinition(
	connectionType internal.ConnectionType,
	stackCancelFunc internal.CancelFunc,
	connectionManager rxgo.IPublishToConnectionManager,
	opts ...rxgo.Option) (*internal.StackDefinition, error) {
	if stackCancelFunc == nil {
		return nil, goerrors.InvalidParam
	}

	nextOutBoundPath := func(ctx context.Context, nextOutboundChannel chan rxgo.Item) connectionWrapper.ConnWrapperNext {
		return func(b []byte) (n int, err error) {
			if ctx.Err() != nil {
				return 0, ctx.Err()
			}
			dataToConnection := gomessageblock.NewReaderWriterSize(len(b))
			if ctx.Err() != nil {
				return 0, ctx.Err()
			}
			n, err = dataToConnection.Write(b)
			if err != nil {
				return 0, err
			}
			item := rxgo.Of(dataToConnection)
			if ctx.Err() != nil {
				return 0, ctx.Err()
			}
			item.SendContext(ctx, nextOutboundChannel)
			return n, nil
		}
	}

	sendPing := func(ctx context.Context, nextOutboundChannel chan rxgo.Item) {
		msg := &wsmsg.WebSocketMessage{
			OpCode: wsmsg.WebSocketMessage_OpPing,
		}
		marshall, err := stream.Marshall(msg)
		if err != nil {
			return
		}
		item := rxgo.Of(marshall)
		item.SendContext(ctx, nextOutboundChannel)
	}
	LastPongReceived := time.Now()
	triggerPingLoop := func(cancelContext context.Context, nextOutboundChannel chan rxgo.Item) {
		ticker := time.NewTicker(time.Second * 5)
		defer ticker.Stop()
		for true {
			select {
			case <-cancelContext.Done():
				return
			case <-ticker.C:
				if time.Now().Sub(LastPongReceived) > time.Second*10 {
					stackCancelFunc("pong time out", true, goerrors.TimeOut)
					return
				}
				sendPing(cancelContext, nextOutboundChannel)
			}
		}
	}

	connectionLoop := func(conn net.Conn, ctx context.Context, nextInBoundChannel chan rxgo.Item) {
		sendMessage := func(message *wsmsg.WebSocketMessage) {
			stm, err := stream.Marshall(message)
			if err != nil {
				//return
			}
			item := rxgo.Of(stm)
			item.SendContext(ctx, nextInBoundChannel)
		}
		message := wsmsg.WebSocketMessage{
			OpCode:  wsmsg.WebSocketMessage_OpStartLoop,
			Message: nil,
		}
		sendMessage(&message)

		for {
			msgs, err := wsutil.ReadServerMessage(conn, nil)
			if err != nil {
				return
			}
			if ctx.Err() != nil {
				return
			}
			for _, msg := range msgs {
				if ctx.Err() != nil {
					return
				}
				var message wsmsg.WebSocketMessage
				switch msg.OpCode {
				case ws.OpContinuation:
					message = wsmsg.WebSocketMessage{
						OpCode:  wsmsg.WebSocketMessage_OpContinuation,
						Message: msg.Payload,
					}
				case ws.OpText:
					message = wsmsg.WebSocketMessage{
						OpCode:  wsmsg.WebSocketMessage_OpText,
						Message: msg.Payload,
					}
				case ws.OpBinary:
					message = wsmsg.WebSocketMessage{
						OpCode:  wsmsg.WebSocketMessage_OpBinary,
						Message: msg.Payload,
					}
				case ws.OpClose:
					stackCancelFunc(
						"close message received",
						true,
						fmt.Errorf(string(msg.Payload)))
				case ws.OpPing:
					print("+")
					err := wsutil.WriteClientMessage(conn, ws.OpPong, msg.Payload)
					if err != nil {
						stackCancelFunc("creating pong payload", true, err)
						return
					}
					continue
				case ws.OpPong:
					_ = LastPongReceived.UnmarshalBinary(msg.Payload)
					continue
				default:
					continue
				}
				if ctx.Err() != nil {
					return
				}
				sendMessage(&message)
				message.Reset()
			}
		}
	}

	// globals
	var connWrapper *connectionWrapper.ConnWrapper
	var pipeWriteClose io.WriteCloser
	var upgradedConnection net.Conn
	var nextInBoundChannel, nextOutboundChannel, tempStep chan rxgo.Item
	const stackName = "WebSocket"
	return &internal.StackDefinition{
		Name: stackName,
		Inbound: func(index int, ctx context.Context) internal.BoundDefinition {
			nextInBoundChannel = make(chan rxgo.Item)
			return internal.BoundDefinition{
				PipeDefinition: func(params internal.PipeDefinitionParams) (rxgo.Observable, error) {
					if stackCancelFunc == nil {
						return nil, goerrors.InvalidParam
					}
					_ = params.Obs.(rxgo.InOutBoundObservable).DoOnNextInOutBound(
						index,
						params.ConnectionId,
						stackName,
						rxgo.StreamDirectionInbound,
						connectionManager,
						func(ctx context.Context, rws goprotoextra.ReadWriterSize) {
							_, err := io.Copy(pipeWriteClose, rws)
							if err != nil {
								return
							}
						}, opts...)
					nextObs := rxgo.FromChannel(nextInBoundChannel)
					return nextObs, nil
				},
				PipeState: internal.PipeState{
					Start: func(ctx context.Context) error {
						return ctx.Err()
					},
					End: func() error {
						close(nextInBoundChannel)
						return nil
					},
				},
			}
		},
		Outbound: func(index int, ctx context.Context) internal.BoundDefinition {
			nextOutboundChannel = make(chan rxgo.Item)
			tempStep = make(chan rxgo.Item)
			return internal.BoundDefinition{
				PipeDefinition: func(params internal.PipeDefinitionParams) (rxgo.Observable, error) {
					if stackCancelFunc == nil {
						return nil, goerrors.InvalidParam
					}
					_ = params.Obs.(rxgo.InOutBoundObservable).DoOnNextInOutBound(
						index-1,
						params.ConnectionId,
						stackName+"FromUser",
						rxgo.StreamDirectionOutbound,
						connectionManager,
						func(ctx context.Context, size goprotoextra.ReadWriterSize) {
							item := rxgo.Of(size)
							item.SendContext(ctx, tempStep)
						},
						opts...)

					tempObs := rxgo.FromChannel(tempStep)
					_ = tempObs.(rxgo.InOutBoundObservable).DoOnNextInOutBound(
						index,
						params.ConnectionId,
						stackName,
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
										err := wsutil.WriteClientMessage(upgradedConnection, ws.OpText, vv.Data.Message)
										if err != nil {
											stackCancelFunc("creating text payload", false, err)
											return
										}
									case wsmsg.WebSocketMessage_OpPing:
										binary, _ := time.Now().MarshalBinary()
										err := wsutil.WriteClientMessage(upgradedConnection, ws.OpPing, binary)
										if err != nil {
											stackCancelFunc("creating ping payload", false, err)
											return
										}
									}
								}
							}
						}, opts...)
					nextObs := rxgo.FromChannel(nextOutboundChannel)
					return nextObs, nil
				},
				PipeState: internal.PipeState{
					Start: func(ctx context.Context) error {
						return ctx.Err()
					},
					End: func() error {
						close(nextOutboundChannel)
						close(tempStep)
						return nil
					},
				},
			}
		},
		StackState: internal.StackState{
			Start: func(conn net.Conn, url *url.URL, ctx context.Context, cancelFunc internal.CancelFunc) (net.Conn, error) {
				if ctx.Err() != nil {
					return nil, ctx.Err()
				}
				var pipeRead io.Reader
				pipeRead, pipeWriteClose = internal.Pipe(ctx)
				if ctx.Err() != nil {
					return nil, ctx.Err()
				}
				connWrapper = connectionWrapper.NewConnWrapper(
					conn,
					ctx,
					pipeRead,
					nextOutBoundPath(ctx, nextOutboundChannel))
				if ctx.Err() != nil {
					return nil, ctx.Err()
				}
				dialer := ws.Dialer{
					NetDial: func(ctx context.Context, network, addr string) (net.Conn, error) {
						if ctx.Err() != nil {
							return nil, ctx.Err()
						}
						return connWrapper, nil
					},
				}
				if ctx.Err() != nil {
					return nil, ctx.Err()
				}
				var err error
				upgradedConnection, _, _, err = dialer.Dial(ctx, url.String())
				if err != nil {
					return nil, err
				}

				if ctx.Err() != nil {
					return nil, ctx.Err()
				}

				go connectionLoop(upgradedConnection, ctx, nextInBoundChannel)
				go triggerPingLoop(ctx, tempStep)
				return upgradedConnection, ctx.Err()
			},
			End: func() error {
				err := pipeWriteClose.Close()
				err = multierr.Append(err, upgradedConnection.Close())
				err = multierr.Append(err, connWrapper.Close())
				return err
			},
		},
	}, nil
}
