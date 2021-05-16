package internal

import (
	"context"
	"github.com/bhbosman/gocomms/connectionManager"
	"github.com/bhbosman/gocomms/internal"
	"github.com/bhbosman/goprotoextra"
	"github.com/reactivex/rxgo/v2"
	"go.uber.org/fx"
	"math"
)

func InvokeInboundTransportLayer(
	params struct {
		fx.In
		Lifecycle         fx.Lifecycle
		IncomingObs       *internal.IncomingObs
		ChannelManager    *internal.ChannelManager
		CancelFunc        context.CancelFunc
		CancelCtx         context.Context
		ConnectionId      string `name:"ConnectionId"`
		ConnectionManager connectionManager.IConnectionManager
	}) {
	params.Lifecycle.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			if params.CancelCtx.Err() != nil {
				return params.CancelCtx.Err()
			}
			params.IncomingObs.InboundObservable.(rxgo.InOutBoundObservable).DoOnNextInOutBound(
				math.MinInt32,
				params.ConnectionId,
				"InboundTransportLayer",
				rxgo.StreamDirectionInbound,
				params.ConnectionManager,
				func(context context.Context, i goprotoextra.ReadWriterSize) {
					if context.Err() != nil {
						return
					}
					params.ChannelManager.Send(context, rxgo.NewNextExternal(true, i))
				})
			return nil
		},
		OnStop: func(ctx context.Context) error {
			return nil
		},
	})
}
