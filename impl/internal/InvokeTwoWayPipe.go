package internal

import (
	"context"
	"github.com/bhbosman/gocomms/internal"
	"go.uber.org/fx"
)

func InvokeFxLifeCycleOutgoingObsStartStop(
	params struct {
		fx.In
		Lifecycle   fx.Lifecycle
		CancelCtx   context.Context
		OutgoingObs *internal.OutgoingObs
	}) error {
	params.Lifecycle.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			return params.CancelCtx.Err()
		},
		OnStop: func(ctx context.Context) error {
			return params.OutgoingObs.Close()
		},
	})
	return nil
}

func InvokeFxLifeCycleIncomingObsStartStop(
	params struct {
		fx.In
		Lifecycle   fx.Lifecycle
		CancelCtx   context.Context
		IncomingObs *internal.IncomingObs
	}) error {
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
