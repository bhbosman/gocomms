package internal

import (
	"context"
	"github.com/bhbosman/gocomms/internal"
	"go.uber.org/fx"
)

func InvokeChannel(
	params struct {
		fx.In
		Channel    *internal.ChannelManager
		LifeCycle  fx.Lifecycle
		CancelFunc context.CancelFunc
		CancelCtx  context.Context
	}) {

	params.LifeCycle.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			return params.CancelCtx.Err()
		},
		OnStop: func(ctx context.Context) error {
			return params.Channel.Close()
		},
	})

}
