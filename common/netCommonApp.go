package common

import (
	"github.com/bhbosman/goCommsDefinitions"
	"github.com/bhbosman/goConnectionManager"
	"github.com/bhbosman/gocommon/GoFunctionCounter"
	"github.com/bhbosman/gocommon/Services/interfaces"
	fx2 "github.com/bhbosman/gocommon/fx"
	"go.uber.org/fx"
	"go.uber.org/fx/fxevent"
	"go.uber.org/zap"
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
		goCommsDefinitions.ProvideCancelContext(params.ParentContext),
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
