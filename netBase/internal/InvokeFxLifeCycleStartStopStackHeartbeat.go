package internal

import (
	"fmt"
	"github.com/bhbosman/gocommon/GoFunctionCounter"
	"github.com/bhbosman/gocommon/model"
	"github.com/reactivex/rxgo/v2"
	"go.uber.org/fx"
	"golang.org/x/net/context"
	"time"
)

func InvokeFxLifeCycleStartStopStackHeartbeat() fx.Option {
	return fx.Invoke(
		func(
			params struct {
				fx.In
				CancelCtx               context.Context
				Lifecycle               fx.Lifecycle
				ConnectionId            string        `name:"ConnectionId"`
				NextFuncOutBoundChannel rxgo.NextFunc `name:"OutBoundChannel"`
				NextFuncInBoundChannel  rxgo.NextFunc `name:"InBoundChannel"`
				GoFunctionCounter       GoFunctionCounter.IService
			},
		) {
			params.Lifecycle.Append(
				fx.Hook{
					OnStart: func(_ context.Context) error {
						if params.CancelCtx.Err() != nil {
							return params.CancelCtx.Err()
						}

						// this function is part of the GoFunctionCounter count
						go func(
							ctx context.Context,
							NextFuncOutBoundChannel rxgo.NextFunc,
							NextFuncInBoundChannel rxgo.NextFunc,
						) {
							functionName := params.GoFunctionCounter.CreateFunctionName(fmt.Sprintf("Heartbeat for connection %v", params.ConnectionId))
							defer func(GoFunctionCounter GoFunctionCounter.IService, name string) {
								_ = GoFunctionCounter.Remove(name)
							}(params.GoFunctionCounter, functionName)
							_ = params.GoFunctionCounter.Add(functionName)
							//
							ticker := time.NewTicker(time.Second)
							defer ticker.Stop()
						loop:
							for {
								select {
								case <-ctx.Done():
									break loop
								case _, ok := <-ticker.C:
									if !ok {
										break loop
									}

									NextFuncInBoundChannel(
										model.NewPublishRxHandlerCounters(
											params.ConnectionId,
											model.StreamDirectionInbound,
										),
									)

									NextFuncOutBoundChannel(
										model.NewPublishRxHandlerCounters(
											params.ConnectionId,
											model.StreamDirectionOutbound,
										),
									)
								}
							}
						}(params.CancelCtx, params.NextFuncOutBoundChannel, params.NextFuncInBoundChannel)
						return nil
					},
					OnStop: func(_ context.Context) error {
						return nil
					},
				},
			)
		},
	)
}
