package internal

import (
	"github.com/bhbosman/goConnectionManager"
	"github.com/bhbosman/gocommon/GoFunctionCounter"
	"github.com/bhbosman/gocommon/model"
	"github.com/bhbosman/gocommon/rxOverride"
	"github.com/bhbosman/gocomms/RxHandlers"
	"github.com/bhbosman/gocomms/common"
	"github.com/reactivex/rxgo/v2"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"golang.org/x/net/context"
	"net"
)

func ProvideRxNextHandlerForNetConnWrite2(name string) fx.Option {
	return fx.Provide(
		fx.Annotated{
			Name: name,
			Target: func(
				params struct {
					fx.In
					Conn                         net.Conn
					OutgoingObs                  *common.OutgoingObs
					RxOptions                    []rxgo.Option
					PublishConnectionInformation goConnectionManager.IPublishConnectionInformation
					ConnectionCancelFunc         model.ConnectionCancelFunc
					Logger                       *zap.Logger
					Ctx                          context.Context
					ConnectionId                 string `name:"ConnectionId"`
					GoFunctionCounter            GoFunctionCounter.IService
				},
			) (*RxHandlers.RxNextHandler, *common.InvokeWriterHandler, error) {
				var handler *common.InvokeWriterHandler
				var rxHandler *RxHandlers.RxNextHandler
				var err error

				handler = common.NewInvokeOutBoundTransportLayerHandler(params.Conn, params.PublishConnectionInformation)
				rxHandler, err = RxHandlers.NewRxNextHandler(
					"net.conn.write",
					params.ConnectionCancelFunc,
					handler,
					handler.SendData,
					handler.SendError,
					handler.Complete,
					params.Logger)
				if err != nil {
					return nil, nil, err
				}

				rxOverride.ForEach(
					"net.conn.write",
					model.StreamDirectionUnknown,
					params.OutgoingObs.OutboundObservable,
					params.Ctx,
					params.GoFunctionCounter,
					rxHandler.OnSendData,
					rxHandler.OnError,
					rxHandler.OnComplete,
					false,
					params.RxOptions...,
				)
				return rxHandler, handler, nil
			},
		},
	)
}
