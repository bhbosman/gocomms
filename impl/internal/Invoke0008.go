package internal

import (
	"context"
	"github.com/bhbosman/gocomms/intf"
	"go.uber.org/fx"
)

func Invoke0008(
	params struct {
		fx.In
		Lifecycle     fx.Lifecycle
		ClientContext intf.IConnectionReactor
		CancelFunc    context.CancelFunc
		CancelCtx     context.Context
	}) {
	params.Lifecycle.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			if params.CancelCtx.Err() != nil {
				return params.CancelCtx.Err()
			}
			return params.ClientContext.Open()
		},
		OnStop: func(ctx context.Context) error {
			return params.ClientContext.Close()
		},
	})
}
