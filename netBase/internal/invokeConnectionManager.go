package internal

import (
	"github.com/bhbosman/goConnectionManager"
	"github.com/reactivex/rxgo/v2"
	"go.uber.org/fx"
	"golang.org/x/net/context"
)

func InvokeConnectionManager() fx.Option {
	return fx.Invoke(
		func(
			params struct {
				fx.In
				ConnectionManager       goConnectionManager.IService
				CancelFunction          context.CancelFunc
				CancelCtx               context.Context
				ConnectionId            string `name:"ConnectionId"`
				Lifecycle               fx.Lifecycle
				NextFuncOutBoundChannel rxgo.NextFunc `name:"OutBoundChannel"`
				NextFuncInBoundChannel  rxgo.NextFunc `name:"InBoundChannel"`
			},
		) {
			params.Lifecycle.Append(
				fx.Hook{
					OnStart: func(_ context.Context) error {
						return params.ConnectionManager.RegisterConnection(
							params.ConnectionId,
							params.CancelFunction,
							params.CancelCtx,
							params.NextFuncOutBoundChannel,
							params.NextFuncInBoundChannel,
						)
					},
					OnStop: func(_ context.Context) error {
						return params.ConnectionManager.DeregisterConnection(params.ConnectionId)
					},
				},
			)
		},
	)
}
