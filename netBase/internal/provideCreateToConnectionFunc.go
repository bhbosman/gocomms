package internal

import (
	"context"
	"github.com/reactivex/rxgo/v2"
	"go.uber.org/fx"
)

func ProvideCreateToConnectionFunc(name string) fx.Option {
	return fx.Provide(
		fx.Annotated{
			Name: name,
			Target: func(
				params struct {
					fx.In
					CancelCtx    context.Context
					NextFunc     rxgo.NextFunc      `name:"OutBoundChannel"`
					ErrorFunc    rxgo.ErrFunc       `name:"OutBoundChannel"`
					CompleteFunc rxgo.CompletedFunc `name:"OutBoundChannel"`
				},
			) rxgo.NextFunc {
				return func(i interface{}) {
					params.NextFunc(i)
				}
			},
		},
	)
}
