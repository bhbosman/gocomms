package internal

import (
	"github.com/bhbosman/gocommon/GoFunctionCounter"
	"github.com/bhbosman/gocommon/model"
	"github.com/bhbosman/gocomms/RxHandlers"
	"github.com/bhbosman/gocomms/common"
	"go.uber.org/fx"
	"golang.org/x/net/context"
	"net"
)

func InvokeFxLifeCycleReadDataFromConnectionStartStop() fx.Option {
	return fx.Invoke(
		func(
			params struct {
				fx.In
				Conn                        net.Conn
				Lifecycle                   fx.Lifecycle
				ConnectionCancelFunc        model.ConnectionCancelFunc
				CancelCtx                   context.Context
				RxNextHandlerForNetConnRead *RxHandlers.RxNextHandler `name:"net.conn.read"`
				GoFunctionCounter           GoFunctionCounter.IService
			},
		) {
			params.Lifecycle.Append(
				fx.Hook{
					OnStart: func(_ context.Context) error {
						if params.CancelCtx.Err() != nil {
							return params.CancelCtx.Err()
						}

						// this function is part of the GoFunctionCounter count
						go func() {
							functionName := params.GoFunctionCounter.CreateFunctionName("InvokeFxLifeCycleReadDataFromConnectionStartStop.ReadFromIoReader")
							defer func(GoFunctionCounter GoFunctionCounter.IService, name string) {
								_ = GoFunctionCounter.Remove(name)
							}(params.GoFunctionCounter, functionName)
							_ = params.GoFunctionCounter.Add(functionName)

							//
							common.ReadFromIoReader(
								"net.conn.read",
								params.Conn,
								params.CancelCtx,
								params.ConnectionCancelFunc,
								params.RxNextHandlerForNetConnRead,
							)

						}()

						return nil
					},
					OnStop: func(_ context.Context) error {
						return params.Conn.Close()
					},
				},
			)
		},
	)
}
