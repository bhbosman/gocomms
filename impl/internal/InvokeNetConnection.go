package internal

import (
	"context"
	"go.uber.org/fx"
	"net"
)

func InvokeNetConnection(params struct {
	fx.In
	Conn       net.Conn
	Lifecycle  fx.Lifecycle
	CancelFunc context.CancelFunc
	CancelCtx  context.Context
}) {
	params.Lifecycle.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			return params.CancelCtx.Err()
		},
		OnStop: func(ctx context.Context) error {
			return params.Conn.Close()
		},
	})
}
