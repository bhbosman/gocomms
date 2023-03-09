package internal

import (
	"github.com/bhbosman/gocomms/common"
	"go.uber.org/fx"
	"golang.org/x/net/context"
)

func InvokeFxLifeCyclePipeStateOutStartStop() fx.Option {
	return fx.Invoke(
		func(
			params struct {
				fx.In
				Lifecycle   fx.Lifecycle
				CancelCtx   context.Context
				PipeState   []*common.PipeState `name:"PipeOutStates"`
				SSStackData map[string]*common.StackDataContainer
			},
		) error {

			for _, pipeStart := range params.PipeState {
				localPipeStart := pipeStart
				params.Lifecycle.Append(
					fx.Hook{
						OnStart: func(_ context.Context) error {
							if params.CancelCtx.Err() != nil {
								return params.CancelCtx.Err()
							}
							var stackData interface{} = nil
							var pipeData interface{} = nil
							if container, ok := params.SSStackData[localPipeStart.ID]; ok {
								stackData = container.StackData
								pipeData = container.OutPipeData
							}
							return localPipeStart.OnStart(stackData, pipeData, params.CancelCtx)
						},
						OnStop: func(_ context.Context) error {
							var stackData interface{} = nil
							var pipeData interface{} = nil
							if container, ok := params.SSStackData[localPipeStart.ID]; ok {
								stackData = container.StackData
								pipeData = container.InPipeData
							}
							return localPipeStart.OnEnd(stackData, pipeData)
						},
					},
				)
			}
			return nil
		},
	)
}
