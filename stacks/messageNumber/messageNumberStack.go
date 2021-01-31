package messageNumber

import (
	"context"
	"encoding/binary"
	"fmt"
	"github.com/bhbosman/gocomms/internal"
	"github.com/bhbosman/goerrors"
	"github.com/bhbosman/gomessageblock"
	"github.com/bhbosman/goprotoextra"
	"github.com/reactivex/rxgo/v2"
)

func StackDefinition(
	userContext interface{},
	stackCancelFunc internal.CancelFunc,
	opts ...rxgo.Option) (*internal.StackDefinition, error) {
	if stackCancelFunc == nil {
		return nil, goerrors.InvalidParam
	}
	const stackName = "MessageNumber"

	return &internal.StackDefinition{
		Name: stackName,
		Inbound: func(inOutBoundParams internal.InOutBoundParams) internal.BoundDefinition {
			return internal.BoundDefinition{
				PipeDefinition: func(pipeParams internal.PipeDefinitionParams) (rxgo.Observable, error) {
					if stackCancelFunc == nil {
						return nil, goerrors.InvalidParam
					}
					errorState := false
					var number uint64 = 0
					return pipeParams.Obs.(rxgo.InOutBoundObservable).MapInOutBound(
						inOutBoundParams.Index,
						pipeParams.ConnectionId,
						stackName,
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
			}
		},
		Outbound: func(inOutBoundParams internal.InOutBoundParams) internal.BoundDefinition {
			return internal.BoundDefinition{
				PipeDefinition: func(pipeParams internal.PipeDefinitionParams) (rxgo.Observable, error) {
					if stackCancelFunc == nil {
						return nil, goerrors.InvalidParam
					}
					errorState := false
					var number uint64 = 0
					return pipeParams.Obs.(rxgo.InOutBoundObservable).MapInOutBound(
						inOutBoundParams.Index,
						pipeParams.ConnectionId,
						stackName,
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
			}
		},
	}, nil
}
