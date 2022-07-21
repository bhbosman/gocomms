package internal

import (
	"github.com/bhbosman/goConnectionManager"
	"go.uber.org/fx"
	"golang.org/x/net/context"
)

func InvokeConnectionManager() fx.Option {
	return fx.Invoke(
		func(
			params struct {
				fx.In
				ConnectionManager goConnectionManager.IService
				CancelFunction    context.CancelFunc
				CancelCtx         context.Context
				ConnectionId      string `name:"ConnectionId"`
				Lifecycle         fx.Lifecycle
			},
		) {
			params.Lifecycle.Append(
				fx.Hook{
					OnStart: func(_ context.Context) error {
						return params.ConnectionManager.RegisterConnection(
							params.ConnectionId,
							params.CancelFunction,
							params.CancelCtx)
					},
					OnStop: func(_ context.Context) error {
						return params.ConnectionManager.DeregisterConnection(params.ConnectionId)
					},
				},
			)
		},
	)
}
