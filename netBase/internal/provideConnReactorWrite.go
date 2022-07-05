package internal

import (
	"context"
	"fmt"
	"github.com/bhbosman/goConnectionManager"
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
					IncomingObs                  *common.IncomingObs
					ChannelManager               chan rxgo.Item `name:"QWERTY"`
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

				createSendData, createSendError, createComplete, err := RxHandlers.All(
					fmt.Sprintf(
						"ProvideConnReactorWrite %v",
						params.ConnectionId),
					model.StreamDirectionUnknown,
					params.ChannelManager,
					params.Logger,
					params.Ctx,
				)
				if err != nil {
					return nil, nil, err
				}

				handler, err = common.NewInvokeInboundTransportLayerHandler(
					params.PublishConnectionInformation,
					createSendData,
					func(i interface{}) error {
						return nil
					},
					createSendError,
					createComplete)
				if err != nil {
					return nil, nil, err
				}

				rxHandler, err = RxHandlers.NewRxNextHandler(
					"conn.reactor.write",
					params.ConnectionCancelFunc,
					handler,
					handler.SendData,
					handler.SendError,
					handler.Complete,
					params.Logger)
				if err != nil {
					return nil, nil, err
				}
				_ = rxOverride.ForEach(
					"conn.reactor.write",
					model.StreamDirectionUnknown,
					params.IncomingObs.InboundObservable,
					params.Ctx,
					params.GoFunctionCounter,
					rxHandler.OnSendData,
					rxHandler.OnError,
					rxHandler.OnComplete,
					false,
					params.RxOptions...)
				return rxHandler, handler, nil
			},
		},
	)
}
