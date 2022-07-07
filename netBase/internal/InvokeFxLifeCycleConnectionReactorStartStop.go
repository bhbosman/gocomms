package internal

import (
	"github.com/bhbosman/gocomms/intf"
	"go.uber.org/fx"
	"golang.org/x/net/context"
)

func InvokeFxLifeCycleConnectionReactorStartStop() fx.Option {
	return fx.Invoke(
		func(
			params struct {
				fx.In
				Lifecycle     fx.Lifecycle
				ClientContext intf.IConnectionReactor
				CancelFunc    context.CancelFunc
				CancelCtx     context.Context
			},
		) {
			params.Lifecycle.Append(
				fx.Hook{
					OnStart: func(_ context.Context) error {
						if params.CancelCtx.Err() != nil {
							return params.CancelCtx.Err()
						}
						return params.ClientContext.Open()
					},
					OnStop: func(_ context.Context) error {
						return params.ClientContext.Close()
					},
				},
			)
		},
	)
}
