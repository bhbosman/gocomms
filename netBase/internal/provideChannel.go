package internal

import (
	"context"
	"fmt"
	"github.com/bhbosman/gocommon/model"
	"github.com/bhbosman/gocomms/RxHandlers"
	"github.com/reactivex/rxgo/v2"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

func ProvideChannel(name string) fx.Option {
	return fx.Provide(
		fx.Annotated{
			Name: name,
			Target: func(
				params struct {
					fx.In
					Lifecycle    fx.Lifecycle
					Logger       *zap.Logger
					Ctx          context.Context
					ConnectionId string `name:"ConnectionId"`
				},
			) (chan rxgo.Item, rxgo.NextFunc, rxgo.ErrFunc, rxgo.CompletedFunc, error) {
				result := make(chan rxgo.Item)
				params.Lifecycle.Append(fx.Hook{
					OnStart: func(_ context.Context) error {
						return nil
					},
					OnStop: func(_ context.Context) error {
						close(result)
						return nil
					},
				})
				handler, err := RxHandlers.All2(
					fmt.Sprintf(
						"ProvideChannel %v",
						params.ConnectionId),
					model.StreamDirectionUnknown,
					result,
					params.Logger,
					params.Ctx,
					true,
				)
				if err != nil {
					return nil, nil, nil, nil, err
				}
				return result, handler.OnSendData, handler.OnError, handler.OnComplete, nil
			},
		},
	)
}
