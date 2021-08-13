package internal

import (
	"context"
	"github.com/bhbosman/gocomms/internal"
	"github.com/google/uuid"
	"go.uber.org/fx"
)

func InvokePipeStateOutStartStop(
	params struct {
		fx.In
		Lifecycle fx.Lifecycle
		CancelCtx context.Context
		PipeState []*internal.PipeState     `name:"PipeOutStates"`
		PipeData  map[uuid.UUID]interface{} `name:"PipeOutData"`
		StackData map[uuid.UUID]interface{} `name:"StackData"`
	}) error {

	for _, pipeStart := range params.PipeState {
		localPipeStart := pipeStart
		params.Lifecycle.Append(fx.Hook{
			OnStart: func(ctx context.Context) error {
				if params.CancelCtx.Err() != nil {
					return params.CancelCtx.Err()
				}
				stackData, _ := params.StackData[localPipeStart.ID]
				pipeData, _ := params.PipeData[localPipeStart.ID]
				return localPipeStart.Start(stackData, pipeData, params.CancelCtx)
			},
			OnStop: func(ctx context.Context) error {
				stackData, _ := params.StackData[localPipeStart.ID]
				pipeData, _ := params.PipeData[localPipeStart.ID]
				return localPipeStart.End(stackData, pipeData)
			},
		})
	}

	return nil
}
