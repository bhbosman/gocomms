package messageNumber

import (
	"context"
	"encoding/binary"
	"fmt"
	"github.com/bhbosman/gocomms/internal"
	"github.com/bhbosman/goerrors"
	"github.com/bhbosman/gomessageblock"
	"github.com/bhbosman/goprotoextra"
	"github.com/google/uuid"
	"github.com/reactivex/rxgo/v2"
)

const StackName = "MessageNumber"

type StackData struct {
}

func StackDefinition(
	connectionType internal.ConnectionType,
	userContext interface{},
	stackCancelFunc internal.CancelFunc,
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
							errorState := false
							var number uint64 = 0
							return id,
								pipeParams.Obs.(rxgo.InOutBoundObservable).MapInOutBound(
									inOutBoundParams.Index,
									pipeParams.ConnectionId,
									StackName,
									rxgo.StreamDirectionInbound,
									pipeParams.ConnectionManager,
									func(ctx context.Context, i goprotoextra.ReadWriterSize) (goprotoextra.ReadWriterSize, error) {
										if errorState {
											stackCancelFunc("In error state", true, goerrors.InvalidState)
											return nil, goerrors.InvalidState
										}
										buffer := [8]byte{0, 0, 0, 0, 0, 0, 0, 0}
										_, err := i.Read(buffer[:])
										if err != nil {
											stackCancelFunc("Could not read message number", true, err)
											errorState = true
											return nil, err
										}
										newNumber := binary.LittleEndian.Uint64(buffer[:])
										number++
										if newNumber != number {
											stackCancelFunc(
												fmt.Sprintf("Invalid number. Expected: %v, Received: %v", number, newNumber),
												true,
												err)
											errorState = true
											return nil, goerrors.InvalidSequenceNumber
										}
										return i, nil
									},
									opts...), nil
						},
					},
					nil
			}),
			Outbound: internal.NewBoundResultImpl(func(inOutBoundParams internal.InOutBoundParams) (internal.IStackBoundDefinition, error) {
				return &internal.StackBoundDefinition{
						PipeDefinition: func(stackData, pipeData interface{}, pipeParams internal.PipeDefinitionParams) (uuid.UUID, rxgo.Observable, error) {
							errorState := false
							var number uint64 = 0
							return id,
								pipeParams.Obs.(rxgo.InOutBoundObservable).MapInOutBound(
									inOutBoundParams.Index,
									pipeParams.ConnectionId,
									StackName,
									rxgo.StreamDirectionOutbound,
									pipeParams.ConnectionManager,
									func(ctx context.Context, i goprotoextra.ReadWriterSize) (goprotoextra.ReadWriterSize, error) {
										if errorState {
											stackCancelFunc("In error state", true, goerrors.InvalidState)
											return nil, goerrors.InvalidState
										}
										number++
										buffer := [8]byte{}
										binary.LittleEndian.PutUint64(buffer[:], number)
										rw := gomessageblock.NewReaderWriterBlock(buffer[:])
										_ = rw.SetNext(i)
										return rw, nil
									},
									opts...), nil
						},
					},
					nil
			}),
		},
		nil
}
