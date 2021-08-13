package pingPong

import (
	"context"
	pingpong "github.com/bhbosman/goMessages/pingpong/stream"
	"github.com/bhbosman/gocommon/stream"
	"github.com/bhbosman/gocomms/internal"
	"github.com/bhbosman/gocomms/intf"
	"github.com/bhbosman/goerrors"
	"github.com/bhbosman/gomessageblock"
	"github.com/bhbosman/goprotoextra"
	"github.com/golang/protobuf/ptypes"
	"github.com/google/uuid"
	"github.com/reactivex/rxgo/v2"
	"go.uber.org/multierr"
	"io"
	"net"
	"net/url"
	"time"
)

func StackDefinition(
	connectionType internal.ConnectionType,
	connectionId string,
	stackCancelFunc internal.CancelFunc,
	opts ...rxgo.Option) (*internal.StackDefinition, error) {
	if stackCancelFunc == nil {
		return nil, goerrors.InvalidParam
	}
	started := false
	var requestId int64 = 0
	id := uuid.New()
	return &internal.StackDefinition{
		IId:  id,
		Name: StackName,
		Inbound: internal.NewBoundResultImpl(func(inOutBoundParams internal.InOutBoundParams) (internal.IStackBoundDefinition, error) {
			return &internal.StackBoundDefinition{
				PipeDefinition: func(stackData, pipeData interface{}, pipeParams internal.PipeDefinitionParams) (uuid.UUID, rxgo.Observable, error) {
					stackDataActual, ok := stackData.(*StackData)
					if !ok {
						return id, nil, WrongStackDataError(connectionType, stackData)
					}
					inBoundChannel := internal.NewChannelManager("inbound PingPong", connectionId)
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
											stackDataActual.outboundChannel.Send(ctx, marshall)
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
			}, nil
		}),
		Outbound: internal.NewBoundResultImpl(func(inOutBoundParams internal.InOutBoundParams) (internal.IStackBoundDefinition, error) {
			return &internal.StackBoundDefinition{
					PipeDefinition: func(stackData, pipeData interface{}, pipeParams internal.PipeDefinitionParams) (uuid.UUID, rxgo.Observable, error) {
						stackDataActual, ok := stackData.(*StackData)
						if !ok {
							return id, nil, WrongStackDataError(connectionType, stackData)
						}
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
								stackDataActual.outboundChannel.Send(ctx, size)
							})
						go func() {
							var ticker *time.Ticker
							ticker = time.NewTicker(time.Second)
							defer ticker.Stop()
							for {
								select {
								case <-disposable:
									return
								case <-pipeParams.CancelContext.Done():
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
										stackDataActual.outboundChannel.Send(pipeParams.CancelContext, marshall)
									}
								}
							}
						}()
						return id, rxgo.FromChannel(stackDataActual.outboundChannel.Items, opts...), nil
					},
				},
				nil
		}),
		StackState: &internal.StackState{
			Id: id,
			Create: func(Conn net.Conn, Url *url.URL, Ctx context.Context, CancelFunc internal.CancelFunc, cfr intf.IConnectionReactorFactoryExtractValues) (interface{}, error) {
				return NewStackData(connectionId), nil
			},
			Destroy: func(stackData interface{}) error {
				var err error = nil
				if _, ok := stackData.(*StackData); !ok {
					err = multierr.Append(err, WrongStackDataError(connectionType, stackData))
				}
				if closer, ok := stackData.(io.Closer); ok {
					err = multierr.Append(err, closer.Close())
				}
				return err
			},
			Start: func(stackData interface{}, startParams internal.StackStartStateParams) (net.Conn, error) {
				if _, ok := stackData.(*StackData); !ok {
					return nil, WrongStackDataError(connectionType, stackData)
				}
				return startParams.Conn, nil
			},
			Stop: func(stackData interface{}, endParams internal.StackEndStateParams) error {
				if _, ok := stackData.(*StackData); !ok {
					return WrongStackDataError(connectionType, stackData)
				}
				return nil
			},
		},
	}, nil
}
