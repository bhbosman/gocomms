package internal

import (
	"context"
	"fmt"
	"github.com/bhbosman/goConnectionManager"
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

func ProvideConnReactorWrite2(name string) fx.Option {
	return fx.Provide(
		fx.Annotated{
			Name: name,
			Target: func(
				params struct {
					fx.In
					IncomingObs                  gocommon.IObservable `name:"Inbound"`
					ChannelManager               chan rxgo.Item       `name:"QWERTY"`
					RxOptions                    []rxgo.Option
					PublishConnectionInformation goConnectionManager.IPublishConnectionInformation
					ConnectionCancelFunc         model.ConnectionCancelFunc
					Logger                       *zap.Logger
					Ctx                          context.Context
					ConnectionId                 string `name:"ConnectionId"`
					GoFunctionCounter            GoFunctionCounter.IService
				},
			) (*RxHandlers.RxNextHandler, *common.InvokeInboundTransportLayerHandler, error) {
				var handler *common.InvokeInboundTransportLayerHandler
				var rxHandler *RxHandlers.RxNextHandler

				eventHandler, err := RxHandlers.All2(
					fmt.Sprintf(
						"ProvideConnReactorWrite %v",
						params.ConnectionId),
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
					params.PublishConnectionInformation,
					eventHandler.OnSendData,
					eventHandler.OnTrySendData,
					func(i interface{}) error {
						return nil
					},
					eventHandler.OnError,
					eventHandler.OnComplete,
				)
				if err != nil {
					return nil, nil, err
				}

				rxHandler, err = RxHandlers.NewRxNextHandler2(
					"ConnectionReactor",
					params.ConnectionCancelFunc,
					handler,
					handler,
					params.Logger)
				if err != nil {
					return nil, nil, err
				}
				_ = rxOverride.ForEach2(
					"conn.reactor.write",
					model.StreamDirectionInbound,
					params.IncomingObs,
					params.Ctx,
					params.GoFunctionCounter,
					rxHandler,
					params.RxOptions...)
				return rxHandler, handler, nil
			},
		},
	)
}
