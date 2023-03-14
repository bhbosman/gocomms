package internal

import (
	"context"
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
	"io"
)

func ProvideRxNextHandlerForNetConnWrite2(name string) fx.Option {
	return fx.Provide(
		fx.Annotated{
			Name: name,
			Target: func(
				params struct {
					fx.In
					Writer                       io.Writer            `name:"PrimaryConnection"`
					OutgoingObs                  gocommon.IObservable `name:"Outbound"`
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

				handler = common.NewInvokeOutBoundTransportLayerHandler(
					params.Writer,
					params.PublishConnectionInformation,
				)
				rxHandler, err = RxHandlers.NewRxNextHandler2(
					"net.conn.write",
					params.ConnectionCancelFunc,
					handler,
					handler,
					params.Logger)
				if err != nil {
					return nil, nil, err
				}

				rxOverride.ForEach2(
					"net.conn.write",
					model.StreamDirectionUnknown,
					params.OutgoingObs,
					params.Ctx,
					params.GoFunctionCounter,
					rxHandler,
					params.RxOptions...,
				)
				return rxHandler, handler, nil
			},
		},
	)
}
