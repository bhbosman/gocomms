package common

import (
	"context"
	"github.com/bhbosman/goCommsDefinitions"
	"github.com/bhbosman/goConn"
	"github.com/bhbosman/goConnectionManager"
	"github.com/bhbosman/gocommon/GoFunctionCounter"
	"github.com/bhbosman/gocommon/Services/interfaces"
	fx2 "github.com/bhbosman/gocommon/fx"
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
	additionalFxOptionsForConnectionInstance func() fx.Option,
	option ...fx.Option,
) fx.Option {
	return fx2.NewFxApplicationOptions(
		startTimeOut,
		stopTimeOut,
		fx.Provide(
			fx.Annotated{
				Target: func() (GoFunctionCounter.IService, error) {
					return params.GoFunctionCounter, nil
				},
			},
		),
		fx.Provide(
			func() (*zap.Logger, error) {
				return params.ZapLogger.Named(connectionName), nil
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
		ProvideCancelContext(params.ParentContext),
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
				CancellationContext goConn.ICancellationContext
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

func ProvideCancelContext(cancelContext context.Context) fx.Option {
	return fx.Provide(
		fx.Annotated{
			Target: func(
				params struct {
					fx.In
					Logger         *zap.Logger
					ConnectionName string `name:"ConnectionName"`
				},
			) (context.Context, context.CancelFunc, goConn.ICancellationContext, error) {
				ctx, cancelFunc := context.WithCancel(cancelContext)
				cancelInstance := goConn.NewCancellationContext(params.ConnectionName, cancelFunc, ctx, params.Logger, nil)
				return ctx,
					func() {
						//todo: move out of callback
						cancelInstance.Cancel("Move out of Callback")
					}, cancelInstance, nil
			},
		},
	)
}
