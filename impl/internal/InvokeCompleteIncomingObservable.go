package internal

import (
	"context"
	"github.com/bhbosman/gocomms/connectionManager"
	"github.com/bhbosman/gocomms/intf"
	"github.com/bhbosman/goprotoextra"
	"github.com/reactivex/rxgo/v2"
	"go.uber.org/fx"
	"net"
	"net/url"
)

func InvokeCompleteIncomingObservable(
	params struct {
		fx.In
		Conn              net.Conn
		Url               *url.URL
		Lifecycle         fx.Lifecycle
		ClientContext     intf.IConnectionReactor
		ConnectionId      string `name:"ConnectionId"`
		ConnectionManager connectionManager.IConnectionManager
		ToConnectionFunc  goprotoextra.ToConnectionFunc
		ToReactorFunc     goprotoextra.ToReactorFunc
		Obs               rxgo.Observable
		CancelCtx         context.Context
	}) {
	params.Lifecycle.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			if params.CancelCtx.Err() != nil {
				return params.CancelCtx.Err()
			}
			onData, err := params.ClientContext.Init(
				params.Conn,
				params.Url,
				params.ConnectionId,
				params.ConnectionManager,
				params.ToConnectionFunc,
				params.ToReactorFunc)
			if err != nil {
				return err
			}
			params.Obs.(rxgo.InOutBoundObservable).DoNextExternal(
				-100,
				params.ConnectionId,
				"ConnectionReactor",
				rxgo.StreamDirectionInbound,
				params.ConnectionManager,
				onData)
			return nil
		},
		OnStop: func(ctx context.Context) error {
			return nil
		},
	})
}
