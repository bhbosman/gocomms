package netBase

import (
	"context"
	"github.com/bhbosman/goCommsDefinitions"
	"go.uber.org/fx"
)

func InvokeCancelContextWithRwc() fx.Option {
	return fx.Invoke(
		func(
			params struct {
				fx.In
				Lifecycle           fx.Lifecycle
				CancellationContext goCommsDefinitions.ICancellationContext
			},
		) error {
			params.Lifecycle.Append(
				fx.Hook{
					OnStart: nil,
					OnStop: func(ctx context.Context) error {
						params.CancellationContext.Cancel()
						return nil
					},
				},
			)
			return nil
		},
	)
}
