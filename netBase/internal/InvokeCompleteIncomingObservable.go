package internal

import (
	"github.com/bhbosman/goConnectionManager"
	"github.com/bhbosman/gocommon/GoFunctionCounter"
	"github.com/bhbosman/gocommon/model"
	"github.com/bhbosman/gocommon/rxOverride"
	"github.com/bhbosman/gocomms/intf"
	"github.com/bhbosman/goprotoextra"
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
				ToConnectionFunc  goprotoextra.ToConnectionFunc `name:"ToConnectionFunc"`
				ToReactorFunc     goprotoextra.ToReactorFunc    `name:"ForReactor"`
				Obs               rxgo.Observable               `name:"QWERTY"`
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
							params.ToConnectionFunc,
							params.ToReactorFunc,
							nil,
							nil,
						)
						if err != nil {
							return err
						}

						// place holders

						rxOverride.ForEach(
							"sadsdasdas",
							model.StreamDirectionUnknown,
							params.Obs,
							params.CancelCtx,
							params.GoFunctionCounter,
							NextFunc,
							ErrFunc,
							CompletedFunc,
							true,
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
