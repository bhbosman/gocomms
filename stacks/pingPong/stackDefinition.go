package pingPong

import (
	"context"
	pingpong "github.com/bhbosman/goMessages/pingpong/stream"
	"github.com/bhbosman/gocommon/stream"
	"github.com/bhbosman/gocomms/internal"
	"github.com/bhbosman/gomessageblock"
	"github.com/bhbosman/goprotoextra"
	"github.com/golang/protobuf/ptypes"
	"github.com/google/uuid"
	"github.com/reactivex/rxgo/v2"
	"net"
	"time"
)

func StackDefinition(
	connectionType internal.ConnectionType,
	connectionId string,
	opts ...rxgo.Option) (*internal.StackDefinition, error) {
	const StackName = "PingPong"
	started := false
	var requestId int64 = 0
	outboundChannel := internal.NewChannelManager(make(chan rxgo.Item), "outbound PingPong", connectionId)
	id := uuid.New()
	return &internal.StackDefinition{
		Id:   id,
		Name: StackName,
		Inbound: internal.NewBoundResultImpl(
			func(stackData, pipeData interface{}, inOutBoundParams internal.InOutBoundParams) internal.IStackBoundDefinition {
				return &internal.StackBoundDefinition{
					PipeDefinition: func(pipeParams internal.PipeDefinitionParams) (uuid.UUID, rxgo.Observable, error) {
						inBoundChannel := internal.NewChannelManager(make(chan rxgo.Item), "inbound PingPong", connectionId)
						disposable := pipeParams.Obs.(rxgo.InOutBoundObservable).DoOnNextInOutBound(
							inOutBoundParams.Index,
							pipeParams.ConnectionId,
							StackName,
							rxgo.StreamDirectionInbound,
							pipeParams.ConnectionManager,
							func(ctx context.Context, incoming goprotoextra.ReadWriterSize) {
								if ctx.Err() != nil {
									return
								}
								if rws, ok := incoming.(*gomessageblock.ReaderWriter); ok {
									tc, err := rws.ReadTypeCode()
									if err != nil {
										return
									}
									if started {
										switch tc {
										case pingpong.PingTypeCode, pingpong.PongTypeCode:
											msg, err := stream.UnMarshal(rws, nil, nil, nil, nil)
											if err != nil {
												return
											}
											switch v := msg.(type) {
											case *pingpong.PingWrapper:
												pong := &pingpong.Pong{
													RequestId:         v.Data.RequestId,
													RequestTimeStamp:  v.Data.RequestTimeStamp,
													ResponseTimeStamp: ptypes.TimestampNow(),
												}
												marshall, err := stream.Marshall(pong)
												if err != nil {
													return
												}
												outboundChannel.Send(ctx, marshall)
											}
											return
										}
									}
								}
								inBoundChannel.Send(ctx, incoming)
							}, opts...)
						go func() {
							<-disposable
							_ = inBoundChannel.Close()
						}()
						obs := rxgo.FromChannel(inBoundChannel.Items, opts...)
						return id, obs, nil
					},
				}
			}),
		Outbound: internal.NewBoundResultImpl(
			func(stackData, pipeData interface{}, inOutBoundParams internal.InOutBoundParams) internal.IStackBoundDefinition {
				return &internal.StackBoundDefinition{
					PipeDefinition: func(pipeParams internal.PipeDefinitionParams) (uuid.UUID, rxgo.Observable, error) {
						disposable := pipeParams.Obs.(rxgo.InOutBoundObservable).DoOnNextInOutBound(
							inOutBoundParams.Index,
							pipeParams.ConnectionId,
							StackName,
							rxgo.StreamDirectionOutbound,
							pipeParams.ConnectionManager,
							func(ctx context.Context, size goprotoextra.ReadWriterSize) {
								if ctx.Err() != nil {
									return
								}
								outboundChannel.Send(ctx, size)
							})
						go func() {
							var ticker *time.Ticker
							ticker = time.NewTicker(time.Second)
							defer ticker.Stop()
							for {
								select {
								case <-disposable:
									return
								case <-inOutBoundParams.Context.Done():
									return
								case <-ticker.C:
									if started {
										requestId++
										tempId := requestId
										ping := &pingpong.Ping{
											RequestId:        tempId,
											RequestTimeStamp: ptypes.TimestampNow(),
										}
										marshall, err := stream.Marshall(ping)
										if err != nil {
											continue
										}
										outboundChannel.Send(inOutBoundParams.Context, marshall)
									}
								}
							}
						}()
						return id,
							rxgo.FromChannel(outboundChannel.Items, opts...),
							nil
					},
				}
			}),
		StackState: internal.StackState{
			Start: func(startParams internal.StackStartStateParams) (net.Conn, error) {
				started = true
				return startParams.Conn, startParams.Ctx.Err()
			},
			End: func(endParams internal.StackEndStateParams) error {
				return outboundChannel.Close()
			},
		},
	}, nil
}
