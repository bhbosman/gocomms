package internal

import (
	"github.com/bhbosman/goCommsDefinitions"
	"github.com/bhbosman/goConnectionManager"
	"github.com/bhbosman/gocommon/GoFunctionCounter"
	"github.com/bhbosman/gocommon/model"
	"github.com/bhbosman/gocommon/rxOverride"
	"github.com/bhbosman/gocomms/intf"
	"github.com/reactivex/rxgo/v2"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"golang.org/x/net/context"
	"net/url"
)

func InvokeCompleteIncomingObservable() fx.Option {
	return fx.Invoke(
		func(
			params struct {
				fx.In
				Url               *url.URL
				Lifecycle         fx.Lifecycle
				ClientContext     intf.IConnectionReactor
				ConnectionId      string `name:"ConnectionId"`
				ConnectionManager goConnectionManager.IService
				ToConnectionFunc  rxgo.NextFunc                   `name:"ToConnectionFunc"`
				ToReactorFunc     rxgo.NextFunc                   `name:"ForReactor"`
				Obs               rxgo.Observable                 `name:"QWERTY"`
				TryNextFunc       goCommsDefinitions.TryNextFunc  `name:"QWERTY"`
				IsNextActive      goCommsDefinitions.IsNextActive `name:"QWERTY"`
				CancelCtx         context.Context
				CancelFunc        context.CancelFunc
				Logger            *zap.Logger
				GoFunctionCounter GoFunctionCounter.IService
				RxOptions         []rxgo.Option
			},
		) {
			params.Lifecycle.Append(
				fx.Hook{
					OnStart: func(
						_ context.Context,
					) error {
						if params.CancelCtx.Err() != nil {
							return params.CancelCtx.Err()
						}
						NextFunc, ErrFunc, CompletedFunc, err := params.ClientContext.Init(
							// this is a thread switch to protect the incoming channel at all cost
							// what happens is that the reactor will send messages to itself, and if the channel is full,
							// we will hang. doing a thread switch is an easy way to fix it, but a better mechanism is required
							// TODO: try to figure out how to protect the channel, without doing a thread switch
							// something like
							//for {
							//	select {
							//	case channel that
							//		params.ToReactorFunc
							//		uses
							//	case context
							//	case time out
							//	}
							//}
							func(i interface{}) {
								params.GoFunctionCounter.GoRun(
									"ThreadSwitch to protect incoming channel",
									func() {
										params.ToReactorFunc(i)
									},
								)
							},
							params.ToConnectionFunc,
						)
						if err != nil {
							return err
						}
						handler := goCommsDefinitions.NewDefaultRxNextHandler(
							NextFunc,
							params.TryNextFunc,
							ErrFunc,
							CompletedFunc,
							params.IsNextActive,
						)
						rxOverride.ForEach2(
							"InvokeCompleteIncomingObservable for Connection ID",
							model.StreamDirectionUnknown,
							params.Obs,
							params.CancelCtx,
							params.GoFunctionCounter,
							handler,
							params.RxOptions...,
						)
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
