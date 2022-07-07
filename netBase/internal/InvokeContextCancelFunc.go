package internal

import (
	"go.uber.org/fx"
	"golang.org/x/net/context"
)

func InvokeContextCancelFunc() fx.Option {
	return fx.Invoke(
		func(
			params struct {
				fx.In
				CancelFunc context.CancelFunc
				LifeCycle  fx.Lifecycle
			},
		) error {
			params.LifeCycle.Append(
				fx.Hook{
					OnStart: nil,
					OnStop: func(ctx context.Context) error {
						if params.CancelFunc != nil {
							params.CancelFunc()
						}
						return nil
					},
				},
			)
			return nil
		},
	)
}
