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

func ProvideCreateChannel(name string) fx.Option {
	return fx.Provide(
		fx.Annotated{
			Name: name,
			Target: func(
				params struct {
					fx.In
					CancelCtx    context.Context
					Logger       *zap.Logger
					ConnectionId string `name:"ConnectionId"`
				},
			) (rxgo.Observable, chan rxgo.Item, rxgo.NextFunc, rxgo.ErrFunc, rxgo.CompletedFunc, error) {
				ch := make(chan rxgo.Item)
				obs := rxgo.FromChannel(ch, rxgo.WithContext(params.CancelCtx))

				nextFunc, errFunc, completeFunc, err := RxHandlers.All(
					fmt.Sprintf(
						"ProvideCreateChannel Provider %v",
						params.ConnectionId),
					model.StreamDirectionInbound,
					ch,
					params.Logger,
					params.CancelCtx,
				)
				if err != nil {
					return nil, nil, nil, nil, nil, err
				}
				return obs, ch, nextFunc, errFunc, completeFunc, nil
			},
		},
	)
}
