package internal

import (
	"github.com/bhbosman/gocomms/RxHandlers"
	"github.com/bhbosman/gocomms/common"
	"go.uber.org/fx"
	"golang.org/x/net/context"
)

func InvokeInboundTransportLayer() fx.Option {
	return fx.Invoke(
		func(
			params struct {
				fx.In
				Lifecycle fx.Lifecycle
				CancelCtx context.Context
				Handler   *common.InvokeInboundTransportLayerHandler `name:"conn.reactor.write"`
				RxHandler *RxHandlers.RxNextHandler                  `name:"conn.reactor.write"`
			}) {
			params.Lifecycle.Append(
				fx.Hook{
					OnStart: func(_ context.Context) error {
						if params.CancelCtx.Err() != nil {
							return params.CancelCtx.Err()
						}

						return nil
					},
					OnStop: func(_ context.Context) error {
						return params.RxHandler.Close()
					},
				},
			)
		},
	)
}
