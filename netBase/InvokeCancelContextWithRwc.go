package netBase

import (
	"context"
	"github.com/bhbosman/gocommon"
	"go.uber.org/fx"
)

func InvokeCancelContextWithRwc() fx.Option {
	return fx.Invoke(
		func(
			params struct {
				fx.In
				Lifecycle           fx.Lifecycle
				CancellationContext gocommon.ICancellationContext
			},
		) error {
			params.Lifecycle.Append(
				fx.Hook{
					OnStart: nil,
					OnStop: func(ctx context.Context) error {
						params.CancellationContext.Cancel("InvokeCancelContextWithRwc")
						return nil
					},
				},
			)
			return nil
		},
	)
}
