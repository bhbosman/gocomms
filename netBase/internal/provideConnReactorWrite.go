package internal

import (
	"context"
	"github.com/bhbosman/gocommon"
	"github.com/bhbosman/gocommon/GoFunctionCounter"
	"github.com/bhbosman/gocommon/model"
	"github.com/bhbosman/gocommon/rxOverride"
	"github.com/bhbosman/gocomms/RxHandlers"
	"github.com/bhbosman/gocomms/common"
	"github.com/reactivex/rxgo/v2"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

func ProvideConnReactorWrite2() fx.Option {
	return fx.Provide(
		fx.Annotated{
			Name: "conn.reactor.write",
			Target: func(
				params struct {
					fx.In
					Observable           gocommon.IObservable `name:"Inbound"`
					ChannelManager       chan rxgo.Item       `name:"QWERTY"`
					RxOptions            []rxgo.Option
					ConnectionCancelFunc model.ConnectionCancelFunc
					Logger               *zap.Logger
					Ctx                  context.Context
					ConnectionId         string `name:"ConnectionId"`
					GoFunctionCounter    GoFunctionCounter.IService
				},
			) (*RxHandlers.RxNextHandler, *common.InvokeInboundTransportLayerHandler, error) {
				var handler *common.InvokeInboundTransportLayerHandler
				var rxHandler *RxHandlers.RxNextHandler

				eventHandler, err := RxHandlers.All2(
					"Deprecated",
					model.StreamDirectionUnknown,
					params.ChannelManager,
					params.Logger,
					params.Ctx,
					true,
				)
				if err != nil {
					return nil, nil, err
				}

				handler, err = common.NewInvokeInboundTransportLayerHandler(
					eventHandler.OnSendData,
					eventHandler.OnTrySendData,
					eventHandler.OnError,
					eventHandler.OnComplete,
				)
				if err != nil {
					return nil, nil, err
				}

				rxHandler, err = RxHandlers.NewRxNextHandler2(
					"Deprecated",
					params.ConnectionCancelFunc,
					handler,
					handler,
					params.Logger)
				if err != nil {
					return nil, nil, err
				}
				_ = rxOverride.ForEach2(
					"Deprecated",
					model.StreamDirectionInbound,
					params.Observable,
					params.Ctx,
					params.GoFunctionCounter,
					rxHandler,
					params.RxOptions...)
				return rxHandler, handler, nil
			},
		},
	)
}
