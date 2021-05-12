package bvisMessageBreaker

import (
	"bytes"
	"context"
	"encoding/binary"
	"fmt"
	internal2 "github.com/bhbosman/gocomms/internal"
	"github.com/bhbosman/gocomms/stacks/bvisMessageBreaker/internal"
	"github.com/bhbosman/goerrors"
	"github.com/bhbosman/gomessageblock"
	"github.com/bhbosman/goprotoextra"
	"github.com/google/uuid"
	"github.com/reactivex/rxgo/v2"
	"reflect"
)

func StackDefinition(
	connectionType internal2.ConnectionType,
	connectionId string,
	stackCancelFunc internal2.CancelFunc,
	stateFunc func(stateFrom, stateTo internal.BuildMessageState, length uint32),
	connectionManager rxgo.IPublishToConnectionManager,
	opts ...rxgo.Option) (*internal2.StackDefinition, error) {
	if stackCancelFunc == nil {
		return nil, goerrors.InvalidParam
	}

	marker := [4]byte{'B', 'V', 'I', 'S'}
	markerAsUInt32 := binary.LittleEndian.Uint32(marker[:])
	id := uuid.New()
	return &internal2.StackDefinition{
		Id:   id,
		Name: StackName,
		Inbound: internal2.NewBoundResultImpl(func(stackData, pipeData interface{}, inOutBoundParams internal2.InOutBoundParams) internal2.IStackBoundDefinition {
			channelManager := internal2.NewChannelManager(make(chan rxgo.Item), "Inbound MessageBreaker", connectionId)
			return &internal2.StackBoundDefinition{
				PipeDefinition: func(params internal2.PipeDefinitionParams) (uuid.UUID, rxgo.Observable, error) {
					if stackCancelFunc == nil {
						return uuid.Nil, nil, goerrors.InvalidParam
					}
					rw := gomessageblock.NewReaderWriter()
					state := internal.BuildMessageStateReadMessageSignature
					var length uint32 = 0

					errorState := false
					inboundState := func(onNext func(data []byte)) {
						var p [4]byte
						canContinue := true
						for canContinue {
							switch state {
							case 0:
								if rw.Size() >= 4 {
									_, err := rw.Read(p[:])
									if err != nil {
										stackCancelFunc("Could not read signature", true, err)
										errorState = true
										return
									}
									c := bytes.Compare(p[:], marker[:])
									if c != 0 {
										stackCancelFunc("Signature incorrect", true, goerrors.InvalidSignature)
										errorState = true
										return
									}
									prev := state
									state = internal.BuildMessageStateReadMessageLength
									if stateFunc != nil {
										stateFunc(prev, state, length)
									}
									break
								} else {
									canContinue = false
								}
							case 1:
								if rw.Size() >= 4 {
									_, err := rw.Read(p[:])
									if err != nil {
										stackCancelFunc("Could not read length", true, err)
										errorState = true
										return
									}
									length = binary.LittleEndian.Uint32(p[:])
									prev := state
									state = internal.BuildMessageStateReadMessageData
									if stateFunc != nil {
										stateFunc(prev, state, length)
									}
									break
								} else {
									canContinue = false
								}
							case 2:
								if uint32(rw.Size()) >= length {
									dataBlock := make([]byte, length)
									_, err := rw.Read(dataBlock)
									if err != nil {
										stackCancelFunc("Could not read data block", true, err)
										errorState = true
										return
									}
									onNext(dataBlock)

									length = 0
									prev := state
									state = internal.BuildMessageStateReadMessageSignature
									if stateFunc != nil {
										stateFunc(prev, state, length)
									}
									break
								} else {
									canContinue = false
								}
							}
						}
					}

					_ = params.Obs.(rxgo.InOutBoundObservable).DoOnNextInOutBound(
						inOutBoundParams.Index,
						params.ConnectionId,
						StackName,
						rxgo.StreamDirectionInbound,
						connectionManager,
						func(ctx context.Context, i goprotoextra.ReadWriterSize) {
							if errorState {
								stackCancelFunc("In error state", true, goerrors.InvalidState)
								return
							}
							switch v := i.(type) {
							case *gomessageblock.ReaderWriter:
								err := rw.SetNext(v)
								if err != nil {
									params.StackCancelFunc("rw.SetNext()", true, err)
									return
								}
								inboundState(func(dataBlock []byte) {
									channelManager.Send(ctx, gomessageblock.NewReaderWriterBlock(dataBlock))
								})
							default:
								stackCancelFunc(
									fmt.Sprintf("Invalid type(%v) received", reflect.TypeOf(i).String()),
									true,
									goerrors.InvalidType)
								errorState = true
								return
							}
						},
						opts...)
					return id,
						rxgo.FromChannel(channelManager.Items, opts...), nil
				},
				PipeState: &internal2.PipeState{
					Start: func(ctx context.Context) error {
						return ctx.Err()
					},
					End: func() error {
						return channelManager.Close()
					},
				},
			}
		}),
		Outbound: internal2.NewBoundResultImpl(
			func(stackData, pipeData interface{}, inOutBoundParams internal2.InOutBoundParams) internal2.IStackBoundDefinition {
				return &internal2.StackBoundDefinition{
					PipeDefinition: func(params internal2.PipeDefinitionParams) (uuid.UUID, rxgo.Observable, error) {
						if stackCancelFunc == nil {
							return uuid.Nil, nil, goerrors.InvalidParam
						}
						errorState := false
						return id,
							params.Obs.(rxgo.InOutBoundObservable).MapInOutBound(
								inOutBoundParams.Index,
								params.ConnectionId,
								StackName,
								rxgo.StreamDirectionOutbound,
								params.ConnectionManager,
								func(ctx context.Context, i goprotoextra.ReadWriterSize) (goprotoextra.ReadWriterSize, error) {
									if errorState {
										stackCancelFunc("In error state", false, goerrors.InvalidState)
										return nil, goerrors.InvalidState
									}
									block := make([]byte, 8)
									binary.LittleEndian.PutUint32(block[0:4], markerAsUInt32)
									binary.LittleEndian.PutUint32(block[4:8], uint32(i.Size()))
									result := gomessageblock.NewReaderWriterBlock(block)
									err := result.SetNext(i)
									if err != nil {
										params.StackCancelFunc("rw.SetNext()", false, err)
										return nil, err
									}
									return result, nil
								},
								opts...),
							nil
					},
				}
			}),
	}, nil
}
