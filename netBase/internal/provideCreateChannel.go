package internal

import (
	"context"
	"fmt"
	"github.com/bhbosman/goCommsDefinitions"
	"github.com/bhbosman/gocommon/model"
	"github.com/bhbosman/gocomms/RxHandlers"
	"github.com/reactivex/rxgo/v2"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

func ProvideCreateChannel() fx.Option {
	return fx.Provide(
		fx.Annotated{
			Name: "QWERTY",
			Target: func(
				params struct {
					fx.In
					CancelCtx    context.Context
					Logger       *zap.Logger
					ConnectionId string `name:"ConnectionId"`
				},
			) (rxgo.Observable,
				chan rxgo.Item,
				rxgo.NextFunc,
				goCommsDefinitions.TryNextFunc,
				rxgo.ErrFunc,
				rxgo.CompletedFunc,
				goCommsDefinitions.IsNextActive,
				error,
			) {
				ch := make(chan rxgo.Item)
				obs := rxgo.FromChannel(ch, rxgo.WithContext(params.CancelCtx))

				ev, err := RxHandlers.All2(
					fmt.Sprintf(
						"ProvideCreateChannel Provider %v",
						params.ConnectionId),
					model.StreamDirectionInbound,
					ch,
					params.Logger,
					params.CancelCtx,
					true,
				)
				if err != nil {
					return nil, nil, nil, nil, nil, nil, nil, err
				}
				return obs, ch, ev.OnSendData, ev.OnTrySendData, ev.OnError, ev.OnComplete, ev.IsActive, nil
			},
		},
	)
}
