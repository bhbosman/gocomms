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
	"golang.org/x/net/context"
)

func InvokeCompleteIncomingObservable() fx.Option {
	return fx.Invoke(
		func(
			params struct {
				fx.In
				Lifecycle               fx.Lifecycle
				ClientContext           intf.IConnectionReactor
				ToConnectionFunc        rxgo.NextFunc                   `name:"ToConnectionFunc"`
				ToReactorFunc           rxgo.NextFunc                   `name:"ForReactor"`
				Observable              rxgo.Observable                 `name:"QWERTY"`
				TryNextFunc             goCommsDefinitions.TryNextFunc  `name:"QWERTY"`
				IsNextActive            goCommsDefinitions.IsNextActive `name:"QWERTY"`
				CancelCtx               context.Context
				GoFunctionCounter       GoFunctionCounter.IService
				RxOptions               []rxgo.Option
				Conn                    goConnectionManager.IPublishConnectionInformation
				NextFuncOutBoundChannel rxgo.NextFunc `name:"OutBoundChannel"`
				NextFuncInBoundChannel  rxgo.NextFunc `name:"InBoundChannel"`
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
						_nextFunc, errFunc, completedFunc, err := params.ClientContext.Init(
							intf.NewInitParams(
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
									if params.CancelCtx.Err() != nil {
										return
									}
									params.GoFunctionCounter.GoRun(
										"ThreadSwitch to protect incoming channel",
										func() {
											if params.CancelCtx.Err() != nil {
												return
											}
											params.ToReactorFunc(i)
										},
									)
								},
								params.ToConnectionFunc,
								params.NextFuncOutBoundChannel,
								params.NextFuncInBoundChannel,
							),
						)
						if err != nil {
							return err
						}
						overrideNextFunc := func(i interface{}) {
							_nextFunc(i)
							switch v := i.(type) {
							case *model.PublishRxHandlerCounters:
								_ = params.Conn.ConnectionInformationReceived(v)
								break
							default:
								break
							}
						}
						handler, err := goCommsDefinitions.NewDefaultRxNextHandler(
							overrideNextFunc,
							params.TryNextFunc,
							errFunc,
							completedFunc,
							params.IsNextActive,
						)
						if err != nil {
							return err
						}

						rxOverride.ForEach2(
							"InvokeCompleteIncomingObservable for Connection ID",
							model.StreamDirectionUnknown,
							params.Observable,
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
