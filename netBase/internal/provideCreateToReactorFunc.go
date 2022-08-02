package internal

import (
	"context"
	"github.com/reactivex/rxgo/v2"
	"go.uber.org/fx"
)

func ProvideCreateToReactorFunc(name string) fx.Option {
	return fx.Provide(
		fx.Annotated{
			Name: name,
			Target: func(
				params struct {
					fx.In
					CancelCtx context.Context
					//GoFunctionCounter GoFunctionCounter.IService
					NextFunc rxgo.NextFunc `name:"QWERTY"`
					//ErrFunc           rxgo.ErrFunc       `name:"QWERTY"`
					//CompletedFunc     rxgo.CompletedFunc `name:"QWERTY"`
				},

			) rxgo.NextFunc {
				return func(i interface{}) {
					params.NextFunc(i)
				}
			},
		},
	)
}
