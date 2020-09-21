package pingPong

import (
	"context"
	pingpong "github.com/bhbosman/goMessages/pingpong/stream"
	"github.com/bhbosman/gocommon/stream"
	"github.com/bhbosman/gocomms/internal"
	"github.com/bhbosman/gomessageblock"
	"github.com/bhbosman/goprotoextra"
	"github.com/golang/protobuf/ptypes"
	"github.com/reactivex/rxgo/v2"
	"net"
	"net/url"
	"time"
)

func StackDefinition(opts ...rxgo.Option) (*internal.StackDefinition, error) {
	const StackName = "PingPong"
	started := false
	var requestId int64 = 0
	outboundChannel := make(chan rxgo.Item)
	return &internal.StackDefinition{
		Name: StackName,
		Inbound: func(index int, ctx context.Context) internal.BoundDefinition {
			return internal.BoundDefinition{
				PipeDefinition: func(params internal.PipeDefinitionParams) (rxgo.Observable, error) {
					var inboundChannel chan rxgo.Item
					inboundChannel = make(chan rxgo.Item)
					disposable := params.Obs.(rxgo.InOutBoundObservable).DoOnNextInOutBound(
						index,
						params.ConnectionId,
						StackName,
						rxgo.StreamDirectionInbound,
						params.ConnectionManager,
						func(ctx context.Context, incoming goprotoextra.ReadWriterSize) {
							if ctx.Err() != nil {
								return
							}
							if rws, ok := incoming.(*gomessageblock.ReaderWriter); ok {
								tc, err := rws.ReadTypeCode()
								if err != nil {
									return
								}
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
										item := rxgo.Of(marshall)
										if ctx.Err() != nil {
											return
										}
										item.SendContext(ctx, outboundChannel)
									case *pingpong.PongWrapper:
									}
									return
								}
							}
							item := rxgo.Of(incoming)
							item.SendContext(ctx, inboundChannel)
						}, opts...)
					go func() {
						defer close(inboundChannel)
						<-disposable
					}()
					obs := rxgo.FromChannel(inboundChannel, opts...)
					return obs, nil
				},
			}
		},
		Outbound: func(index int, ctx context.Context) internal.BoundDefinition {
			return internal.BoundDefinition{
				PipeDefinition: func(params internal.PipeDefinitionParams) (rxgo.Observable, error) {
					disposable := params.Obs.(rxgo.InOutBoundObservable).DoOnNextInOutBound(
						index,
						params.ConnectionId,
						StackName,
						rxgo.StreamDirectionOutbound,
						params.ConnectionManager,
						func(ctx context.Context, size goprotoextra.ReadWriterSize) {
							if ctx.Err() != nil {
								return
							}
							item := rxgo.Of(size)
							item.SendContext(ctx, outboundChannel)
						})
					go func() {
						var ticker *time.Ticker
						ticker = time.NewTicker(time.Second)
						defer ticker.Stop()
						for {
							select {
							case <-disposable:
								return
							case <-ctx.Done():
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
									item := rxgo.Of(marshall)
									item.SendContext(ctx, outboundChannel)
								}
							}
						}
					}()
					return rxgo.FromChannel(outboundChannel, opts...), nil
				},
			}
		},
		StackState: internal.StackState{
			Start: func(conn net.Conn, url *url.URL, ctx context.Context, cancelFunc internal.CancelFunc) (net.Conn, error) {
				started = true
				return conn, ctx.Err()
			},
			End: func() error {
				started = false
				defer close(outboundChannel)
				return nil
			},
		},
	}, nil
}
