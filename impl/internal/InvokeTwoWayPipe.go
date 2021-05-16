package internal

import (
	"context"
	"github.com/bhbosman/gocomms/internal"
	"go.uber.org/fx"
)

func InvokeTwoWayPipe(
	params struct {
		fx.In
		Lifecycle   fx.Lifecycle
		CancelCtx   context.Context
		OutgoingObs *internal.OutgoingObs
		IncomingObs *internal.IncomingObs
	}) error {
	params.Lifecycle.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			return params.CancelCtx.Err()
		},
		OnStop: func(ctx context.Context) error {
			return params.OutgoingObs.Close()
		},
	})
	params.Lifecycle.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			return params.CancelCtx.Err()
		},
		OnStop: func(ctx context.Context) error {
			return params.IncomingObs.Close()
		},
	})
	return nil
}
