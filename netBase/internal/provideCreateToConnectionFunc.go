package internal

import (
	"context"
	"github.com/bhbosman/goprotoextra"
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
			) goprotoextra.ToConnectionFunc {
				return func(rw goprotoextra.ReadWriterSize) error {
					err := params.CancelCtx.Err()
					if err != nil {
						return err
					}
					params.NextFunc(rw)
					return nil
				}
			},
		},
	)
}
