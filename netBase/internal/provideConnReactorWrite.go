package internal

import (
	"context"
	"github.com/bhbosman/gocommon"
	"github.com/bhbosman/gocommon/GoFunctionCounter"
	"github.com/bhbosman/gocommon/model"
	"github.com/bhbosman/gocommon/rxOverride"
	"github.com/bhbosman/gocomms/RxHandlers"
	"github.com/bhbosman/gocomms/common"
	rxgo "github.com/reactivex/rxgo/v2"
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

				eventHandler, err := RxHandlers.All2(
					"ConnReactorWrite",
					model.StreamDirectionUnknown,
					params.ChannelManager,
					params.Logger,
					params.Ctx,
					true,
				)
				if err != nil {
					return nil, nil, err
				}

				handler, err := common.NewInvokeInboundTransportLayerHandler(
					eventHandler.OnSendData,
					eventHandler.OnTrySendData,
					eventHandler.OnError,
					eventHandler.OnComplete,
				)
				if err != nil {
					return nil, nil, err
				}

				rxHandler, err := RxHandlers.NewRxNextHandler2(
					"ConnReactorWrite",
					params.ConnectionCancelFunc,
					handler,
					handler,
					params.Logger)
				if err != nil {
					return nil, nil, err
				}
				_ = rxOverride.ForEach2(
					"ConnReactorWrite",
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
