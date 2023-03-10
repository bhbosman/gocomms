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
			var ticker *time.Ticker
			params.Lifecycle.Append(
				fx.Hook{
					OnStart: func(_ context.Context) error {
						ticker = time.NewTicker(time.Second)
						if params.CancelCtx.Err() != nil {
							return params.CancelCtx.Err()
						}
						return params.GoFunctionCounter.GoRun(
							fmt.Sprintf("Heartbeat for connection %v", params.ConnectionId),
							func() {
							loop:
								for {
									select {
									case <-params.CancelCtx.Done():
										break loop
									case _, ok := <-ticker.C:
										if !ok {
											break loop
										}
										params.NextFuncInBoundChannel(
											model.NewPublishRxHandlerCounters(
												params.ConnectionId,
												model.StreamDirectionInbound,
											),
										)
										params.NextFuncOutBoundChannel(
											model.NewPublishRxHandlerCounters(
												params.ConnectionId,
												model.StreamDirectionOutbound,
											),
										)
									}
								}
							},
						)
					},
					OnStop: func(_ context.Context) error {
						if ticker != nil {
							ticker.Stop()
						}
						return nil
					},
				},
			)
		},
	)
}
