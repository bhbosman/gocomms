package internal

import (
	"fmt"
	"github.com/bhbosman/goCommsDefinitions"
	"github.com/bhbosman/gocommon"
	"github.com/bhbosman/gocommon/model"
	"github.com/bhbosman/gocomms/RxHandlers"
	"github.com/reactivex/rxgo/v2"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"golang.org/x/net/context"
)

func ProvideRxNextHandlerForNetConnRead22(name string) fx.Option {
	return fx.Provide(
		fx.Annotated{
			Name: name,
			Target: func(
				params struct {
					fx.In
					IncomingObs          gocommon.IObservable `name:"Inbound"`
					Ctx                  context.Context
					ConnectionCancelFunc model.ConnectionCancelFunc
					InBoundChannel       chan rxgo.Item `name:"InBoundChannel"`
					Logger               *zap.Logger
					ConnectionId         string `name:"ConnectionId"`
				},
			) (*RxHandlers.RxNextHandler, rxgo.NextFunc, goCommsDefinitions.TryNextFunc, rxgo.ErrFunc, rxgo.CompletedFunc, error) {

				// DO NOT set complete param to RxHandlers.CreateComplete(params.InBoundChannel)
				// as it will lead to a double close and a panic
				// The params.Inbound Channel will eb closed by die FxApp.Shutdown, so no further close() is required
				// make sure useCompleteCallback = false
				eventHandler, err := RxHandlers.All2(
					fmt.Sprintf(
						"ProvideRxNextHandlerForNetConnRead %v",
						params.ConnectionId),
					model.StreamDirectionUnknown,
					params.InBoundChannel,
					params.Logger,
					params.Ctx,
					false,
				)
				if err != nil {
					return nil, nil, nil, nil, nil, err
				}

				result, err := RxHandlers.NewRxNextHandler2(
					"net.conn.read",
					params.ConnectionCancelFunc,
					nil,
					eventHandler, /*see comment*/
					params.Logger)
				if err != nil {
					return nil, nil, nil, nil, nil, err
				}
				return result, eventHandler.OnSendData, eventHandler.OnTrySendData, eventHandler.OnError, eventHandler.OnComplete, nil
			},
		},
	)
}
