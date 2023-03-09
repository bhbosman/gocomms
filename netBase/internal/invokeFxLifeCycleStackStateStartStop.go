package internal

import (
	"fmt"
	"github.com/bhbosman/gocomms/common"
	"github.com/reactivex/rxgo/v2"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"golang.org/x/net/context"
	"net"
)

func InvokeFxLifeCycleStackStateStartStop() fx.Option {
	return fx.Invoke(
		func(
			params struct {
				fx.In
				Lifecycle     fx.Lifecycle
				Conn          net.Conn `name:"PrimaryConnection"`
				CancelCtx     context.Context
				CancelCtxFunc context.CancelFunc
				StackState    []common.IStackState
				StackData     map[string]*common.StackDataContainer
				Logger        *zap.Logger
				ToReactorFunc rxgo.NextFunc `name:"ForReactor"`
			},
		) error {
			localConn := params.Conn

			reportConnectionNilReturned := false
			for _, stackState := range params.StackState {
				if params.CancelCtx.Err() != nil {
					return params.CancelCtx.Err()
				}
				localStackState := stackState
				params.Lifecycle.Append(
					fx.Hook{
						OnStart: func(_ context.Context) error {
							if params.CancelCtx.Err() != nil {
								return params.CancelCtx.Err()
							}
							if localConn == nil {
								return fmt.Errorf("no incoming connection when stsrting stacks")
							}
							if localStackState.OnStart() != nil {
								//var err error
								var stackData common.IStackCreateData
								if container, ok := params.StackData[localStackState.GetId()]; ok {
									stackData = container.StackData
								}
								if localStackState.GetHijackStack() && localConn == nil {
									hijackStackFailedError := fmt.Errorf(
										"stack %v could not start, as it requires a connection to hi-jack and it has been set to nil by a previous stack",
										localStackState.GetId())
									params.Logger.Error(
										"Stack start failed, as it could not hi-jack the connection",
										zap.String("StackName", localStackState.GetId()),
										zap.Error(hijackStackFailedError))
									return hijackStackFailedError
								}
								var err error
								localConn, err = localStackState.OnStart()(localConn, stackData, params.ToReactorFunc)

								if err != nil {
									params.Logger.Error(
										"Stack start failed",
										zap.String("StackName", localStackState.GetId()),
										zap.Error(err))
									return err
								}

								if !reportConnectionNilReturned && localConn == nil {
									reportConnectionNilReturned = true
									params.Logger.Info(
										"Stack state Start method returned a nil localConn. No more hi-jacking of connection allowed",
										zap.String("StackName", localStackState.GetId()))
								}
								return err
							}
							return params.CancelCtx.Err()
						},
						OnStop: func(_ context.Context) error {
							if localStackState.OnStop() != nil {
								var stackData interface{}
								if container, ok := params.StackData[localStackState.GetId()]; ok {
									stackData = container.StackData
								}
								return localStackState.OnStop()(stackData)
							}
							return nil
						},
					},
				)
			}
			return nil
		},
	)
}
