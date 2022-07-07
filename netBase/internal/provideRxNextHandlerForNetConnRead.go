package internal

import (
	"fmt"
	"github.com/bhbosman/gocommon/model"
	"github.com/bhbosman/gocomms/RxHandlers"
	"github.com/bhbosman/gocomms/common"
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
					IncomingObs          *common.IncomingObs
					Ctx                  context.Context
					ConnectionCancelFunc model.ConnectionCancelFunc
					InBoundChannel       chan rxgo.Item `name:"InBoundChannel"`
					Logger               *zap.Logger
					ConnectionId         string `name:"ConnectionId"`
				},
			) (*RxHandlers.RxNextHandler, rxgo.NextFunc, rxgo.ErrFunc, rxgo.CompletedFunc, error) {

				// DO NOT set complete param to RxHandlers.CreateComplete(params.InBoundChannel)
				// as it will lead to a double close and a panic
				// The params.Inbound Channel will eb closed by die FxApp.Shutdown, so no further close() is required

				nextFunc, errFunc, completedFunc, err := RxHandlers.All(
					fmt.Sprintf(
						"ProvideRxNextHandlerForNetConnRead %v",
						params.ConnectionId),
					model.StreamDirectionUnknown,
					params.InBoundChannel,
					params.Logger,
					params.Ctx,
				)
				if err != nil {
					return nil, nil, nil, nil, err
				}

				result, err := RxHandlers.NewRxNextHandler(
					"net.conn.read",
					params.ConnectionCancelFunc,
					nil,
					nextFunc,
					errFunc,
					nil, /*see comment*/
					params.Logger)
				if err != nil {
					return nil, nil, nil, nil, err
				}
				return result, nextFunc, errFunc, completedFunc, nil
			},
		},
	)
}
