package internal

import (
	"context"
	"github.com/bhbosman/gocomms/connectionManager"
	"go.uber.org/fx"
)

func InvokeConnectionManager(
	params struct {
		fx.In
		ConnectionManager connectionManager.IConnectionManager
		CancelFunction    context.CancelFunc
		CancelCtx         context.Context
		ConnectionId      string `name:"ConnectionId"`
		Lifecycle         fx.Lifecycle
	}) {
	params.Lifecycle.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			if params.CancelCtx.Err() != nil {
				return params.CancelCtx.Err()
			}
			return params.ConnectionManager.RegisterConnection(
				params.ConnectionId,
				params.CancelFunction,
				params.CancelCtx)
		},
		OnStop: func(ctx context.Context) error {
			return params.ConnectionManager.DeregisterConnection(params.ConnectionId)
		},
	})
}
