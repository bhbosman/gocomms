package internal

import (
	"context"
	"github.com/bhbosman/gocomms/connectionManager"
	commsInternal "github.com/bhbosman/gocomms/internal"
	"github.com/bhbosman/goprotoextra"
	"github.com/reactivex/rxgo/v2"
	"go.uber.org/fx"
	"io"
	"math"
	"net"
)

func InvokeOutBoundTransportLayer(
	params struct {
		fx.In
		Conn              net.Conn
		Lifecycle         fx.Lifecycle
		CancelFunc        context.CancelFunc
		OutgoingObs       *commsInternal.OutgoingObs
		CancelCtx         context.Context
		ConnectionId      string `name:"ConnectionId"`
		ConnectionManager connectionManager.IConnectionManager
	}) {
	params.Lifecycle.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			if params.CancelCtx.Err() != nil {
				return params.CancelCtx.Err()
			}
			outboundNextClosed := false
			_ = params.OutgoingObs.OutboundObservable.(rxgo.InOutBoundObservable).DoOnNextInOutBound(
				math.MaxInt16,
				params.ConnectionId,
				"Connection Write",
				rxgo.StreamDirectionOutbound,
				params.ConnectionManager,
				func(context context.Context, i goprotoextra.ReadWriterSize) {
					if context.Err() != nil {
						return
					}
					if outboundNextClosed {
						return
					}
					_, err := io.Copy(params.Conn, i)
					if err != nil {
						params.CancelFunc()
						return
					}
				})
			return nil
		},
		OnStop: func(ctx context.Context) error {
			return nil
		},
	})
}
