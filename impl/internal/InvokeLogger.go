package internal

import (
	"context"
	"github.com/bhbosman/gologging"
	"go.uber.org/fx"
)

func InvokeLogger(
	params struct {
		fx.In
		LifeCycle  fx.Lifecycle
		Logger     *gologging.SubSystemLogger
		CancelFunc context.CancelFunc
		CancelCtx  context.Context
	}) {

	params.LifeCycle.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			if params.CancelCtx.Err() != nil {
				return params.CancelCtx.Err()
			}
			return nil
		},
		OnStop: func(ctx context.Context) error {
			return nil
		},
	})
}
