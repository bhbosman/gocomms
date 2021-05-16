package internal

import (
	"context"
	"github.com/bhbosman/gocomms/internal"
	"github.com/bhbosman/gocomms/intf"
	"github.com/google/uuid"
	"go.uber.org/fx"
	"net"
	"net/url"
)

func InvokeStackStateStartStop(
	params struct {
		fx.In
		Url        *url.URL
		Lifecycle  fx.Lifecycle
		Conn       net.Conn
		CancelCtx  context.Context
		CancelFunc internal.CancelFunc
		Cfr        intf.IConnectionReactorFactory
		StackState []*internal.StackState
		StackData  map[uuid.UUID]interface{} `name:"StackData"`
	}) error {
	connectionReactorFactory := params.Cfr
	localConn := params.Conn
	for _, stackState := range params.StackState {
		if params.CancelCtx.Err() != nil {
			return params.CancelCtx.Err()
		}
		localStackState := stackState
		params.Lifecycle.Append(fx.Hook{
			OnStart: func(ctx context.Context) error {
				if params.CancelCtx.Err() != nil {
					return params.CancelCtx.Err()
				}
				if localConn == nil {
					return nil
				}
				if localStackState.Start != nil {
					var err error
					stackData, _ := params.StackData[localStackState.Id]

					localConn, err = localStackState.Start(
						stackData,
						internal.NewStackStartStateParams(
							localConn,
							params.Url,
							params.CancelCtx,
							params.CancelFunc,
							connectionReactorFactory))
					return err
				}
				return params.CancelCtx.Err()
			},
			OnStop: func(ctx context.Context) error {
				if localStackState.Stop != nil {
					stackData, _ := params.StackData[localStackState.Id]
					return localStackState.Stop(stackData, internal.NewStackEndStateParams())
				}
				return nil
			},
		})
	}
	return nil
}
