package internal

import (
	"github.com/bhbosman/gocommon/GoFunctionCounter"
	"github.com/bhbosman/gocommon/model"
	"github.com/bhbosman/gocomms/RxHandlers"
	"github.com/bhbosman/gocomms/common"
	"go.uber.org/fx"
	"golang.org/x/net/context"
	"io"
)

func InvokeFxLifeCycleReadDataFromConnectionStartStop() fx.Option {
	return fx.Invoke(
		func(
			params struct {
				fx.In
				Reader                      io.Reader `name:"PrimaryConnection"`
				Lifecycle                   fx.Lifecycle
				ConnectionCancelFunc        model.ConnectionCancelFunc
				CancelCtx                   context.Context
				CancelFunc                  context.CancelFunc
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
						return params.GoFunctionCounter.GoRun(
							"InvokeFxLifeCycleReadDataFromConnectionStartStop.ReadFromIoReader",
							func() {
								common.ReadFromIoReader(
									params.Reader,
									params.CancelCtx,
									params.CancelFunc,
									params.RxNextHandlerForNetConnRead,
								)
							},
						)
					},
				},
			)
		},
	)
}
