package internal

import (
	"context"
	"github.com/bhbosman/gocommon/GoFunctionCounter"
	"github.com/bhbosman/goprotoextra"
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
					CancelCtx         context.Context
					GoFunctionCounter GoFunctionCounter.IService
					NextFunc          rxgo.NextFunc      `name:"QWERTY"`
					ErrFunc           rxgo.ErrFunc       `name:"QWERTY"`
					CompletedFunc     rxgo.CompletedFunc `name:"QWERTY"`
				},
			) goprotoextra.ToReactorFunc {
				return func(inline bool, any interface{}) error {
					err := params.CancelCtx.Err()
					if err != nil {
						return err
					}
					if !inline {

						// this function is part of the GoFunctionCounter count
						go func(
							any interface{},
						) {
							functionName := params.GoFunctionCounter.CreateFunctionName("ProvideCreateToReactorFunc")
							defer func(GoFunctionCounter GoFunctionCounter.IService, name string) {
								_ = GoFunctionCounter.Remove(name)
							}(params.GoFunctionCounter, functionName)
							_ = params.GoFunctionCounter.Add(functionName)
							params.NextFunc(any)
						}(any)
						return nil
					}
					params.NextFunc(any)
					return nil
				}
			},
		},
	)
}
