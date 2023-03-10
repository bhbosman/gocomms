package internal

import (
	"context"
	"github.com/reactivex/rxgo/v2"
	"go.uber.org/fx"
)

func ProvideCreateToReactorFunc() fx.Option {
	return fx.Provide(
		fx.Annotated{
			Name: "ForReactor",
			Target: func(
				params struct {
					fx.In
					CancelCtx context.Context
					NextFunc  rxgo.NextFunc `name:"QWERTY"`
				},

			) rxgo.NextFunc {
				return func(i interface{}) {
					params.NextFunc(i)
				}
			},
		},
	)
}
