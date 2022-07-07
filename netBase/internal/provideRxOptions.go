package internal

import (
	"context"
	"github.com/reactivex/rxgo/v2"
	"go.uber.org/fx"
)

func ProvideRxOptions() fx.Option {
	return fx.Provide(
		fx.Annotated{
			Target: func(
				params struct {
					fx.In
					Ctx context.Context
				},
			) []rxgo.Option {
				result := []rxgo.Option{
					rxgo.WithContext(params.Ctx),
					rxgo.WithBufferedChannel(1024),
					rxgo.WithBackPressureStrategy(rxgo.Block),
				}
				return result
			},
		},
	)
}
