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
				Url                    *url.URL
				Lifecycle              fx.Lifecycle
				ClientContext          intf.IConnectionReactor
				ConnectionId           string `name:"ConnectionId"`
				ConnectionManager      goConnectionManager.IService
				ToConnectionFunc       rxgo.NextFunc                  `name:"ToConnectionFunc"`
				ToReactorFunc          rxgo.NextFunc                  `name:"ForReactor"`
				Obs                    rxgo.Observable                `name:"QWERTY"`
				NextFuncWithNotSending goCommsDefinitions.TryNextFunc `name:"QWERTY"`
				CancelCtx              context.Context
				CancelFunc             context.CancelFunc
				Logger                 *zap.Logger
				GoFunctionCounter      GoFunctionCounter.IService
				RxOptions              []rxgo.Option
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
						NextFunc, ErrFunc, CompletedFunc, channel, err := params.ClientContext.Init(
							params.ToReactorFunc,
							params.ToConnectionFunc,
						)
						if err != nil {
							return err
						}
						handler := goCommsDefinitions.NewDefaultRxNextHandler(
							NextFunc,
							params.NextFuncWithNotSending,
							ErrFunc,
							CompletedFunc,
							func() bool {
								return true
							},
						)

						// place holders
						if channel == nil {
							rxOverride.ForEach2(
								"InvokeCompleteIncomingObservable for Connection ID",
								model.StreamDirectionUnknown,
								params.Obs,
								params.CancelCtx,
								params.GoFunctionCounter,
								handler,
								params.RxOptions...,
							)

						} else {
							rxOverride.ForEachWithChannel(
								"InvokeCompleteIncomingObservable for Connection ID",
								model.StreamDirectionUnknown,
								params.Obs,
								params.CancelCtx,
								params.GoFunctionCounter,
								NextFunc,
								params.NextFuncWithNotSending,
								ErrFunc,
								CompletedFunc,
								func() bool {
									return true
								},
								channel,
								params.RxOptions...,
							)
						}
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
