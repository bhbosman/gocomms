package internal

import (
	"context"
	"github.com/bhbosman/gocomms/connectionManager"
	"github.com/bhbosman/gocomms/internal"
	"github.com/bhbosman/goprotoextra"
	"go.uber.org/fx"
	"net"
)

func InvokeFxLifeCycleReadDataFromConnectionStartStop(
	params struct {
		fx.In
		Conn              net.Conn
		Lifecycle         fx.Lifecycle
		CancelFunc        internal.CancelFunc
		CancelCtx         context.Context
		IncomingObs       *internal.IncomingObs
		ConnectionManager connectionManager.IConnectionManager
		ConnectionId      string `name:"ConnectionId"`
	}) {
	params.Lifecycle.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			if params.CancelCtx.Err() != nil {
				return params.CancelCtx.Err()
			}
			go internal.ReadDataFromConnection(
				params.Conn,
				params.CancelFunc,
				params.CancelCtx,
				params.ConnectionManager,
				params.ConnectionId,
				1000000,
				"Connection Read",
				func(rws goprotoextra.IReadWriterSize, cancelCtx context.Context, CancelFunc internal.CancelFunc) {
					err := params.IncomingObs.ReceiveIncomingData(rws)
					if err != nil {
						CancelFunc("Connection Read", true, err)
						return
					}
				})
			return nil
		},
		OnStop: func(ctx context.Context) error {
			return params.Conn.Close()
		},
	})
}
