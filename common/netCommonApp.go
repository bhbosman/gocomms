package common

import (
	"context"
	"github.com/bhbosman/goCommsDefinitions"
	"github.com/bhbosman/goConnectionManager"
	"github.com/bhbosman/gocommon"
	"github.com/bhbosman/gocommon/GoFunctionCounter"
	fx2 "github.com/bhbosman/gocommon/fx"
	"github.com/bhbosman/gocommon/services/interfaces"
	"go.uber.org/fx"
	"go.uber.org/fx/fxevent"
	"go.uber.org/zap"
	"net"
	"time"
)

func ConnectionApp(
	startTimeOut time.Duration,
	stopTimeOut time.Duration,
	connectionName string,
	connectionInstancePrefix string,
	params NetAppFuncInParams,
	cancellationContext gocommon.ICancellationContext,
	namedLogger *zap.Logger,
	additionalFxOptionsForConnectionInstance func() fx.Option,
	option ...fx.Option,
) fx.Option {
	return fx2.NewFxApplicationOptions(
		startTimeOut,
		stopTimeOut,
		fx.Provide(
			fx.Annotated{
				Target: func() (context.Context, context.CancelFunc, gocommon.ICancellationContext, error) {
					return cancellationContext.CancelContext(),
						func() {
							//todo: move out of callback
							cancellationContext.Cancel("Move out of Callback")
						},
						cancellationContext, nil
				},
			},
		),
		fx.Provide(
			fx.Annotated{
				Target: func() (GoFunctionCounter.IService, error) {
					return params.GoFunctionCounter, nil
				},
			},
		),
		fx.Provide(
			func() (*zap.Logger, error) {
				return namedLogger, nil
			},
		),
		fx.WithLogger(
			func(logger *zap.Logger) (fxevent.Logger, error) {
				return &fxevent.ZapLogger{
					Logger: logger,
				}, nil
			},
		),
		fx.Provide(
			fx.Annotated{
				Target: func() interfaces.IUniqueReferenceService {
					return params.UniqueSessionNumber
				},
			},
		),
		goCommsDefinitions.ProvideStringContext("ConnectionName", connectionName),
		goCommsDefinitions.ProvideStringContext("ConnectionInstancePrefix", connectionInstancePrefix),
		goConnectionManager.ProvideConnectionManager(params.ConnectionManager),
		goConnectionManager.ProvideCommandsToConnectionManager(),
		goConnectionManager.ProvideObtainConnectionManagerInformation(),
		goConnectionManager.ProvideRegisterToConnectionManager(),
		goConnectionManager.ProvidePublishConnectionInformation(),
		fx.Provide(
			fx.Annotated{
				Target: func() func() fx.Option {
					return additionalFxOptionsForConnectionInstance
				},
			},
		),
		fx.Options(option...),
	)
}

func InvokeCancelContext() fx.Option {
	return fx.Invoke(
		func(
			params struct {
				fx.In
				Lifecycle           fx.Lifecycle
				CancellationContext gocommon.ICancellationContext
			},
		) error {
			params.Lifecycle.Append(
				fx.Hook{
					OnStart: nil,
					OnStop: func(ctx context.Context) error {
						params.CancellationContext.Cancel("InvokeCancelContext")
						return nil
					},
				},
			)
			return nil
		},
	)
}

func InvokeListenerClose() fx.Option {
	return fx.Invoke(
		func(
			params struct {
				fx.In
				NetListener net.Listener
				Lifecycle   fx.Lifecycle
			},
		) {
			params.Lifecycle.Append(
				fx.Hook{
					OnStart: nil,
					OnStop: func(ctx context.Context) error {
						err := params.NetListener.Close()
						return err
					},
				},
			)
		},
	)
}

//func ProvideCancelContext(cancelContext context.Context) fx.Option {
//	return fx.Provide(
//		fx.Annotated{
//			Target: func(
//				params struct {
//				fx.In
//				Logger         *zap.Logger
//				ConnectionName string `name:"ConnectionName"`
//			},
//			) (context.Context, context.CancelFunc, goConn.ICancellationContext, error) {
//				ctx, cancelFunc := context.WithCancel(cancelContext)
//				cancelInstance := goConn.NewCancellationContext(params.ConnectionName, cancelFunc, ctx, params.Logger, nil)
//				return ctx,
//					func() {
//						//todo: move out of callback
//						cancelInstance.Cancel("Move out of Callback")
//					}, cancelInstance, nil
//			},
//		},
//	)
//}
