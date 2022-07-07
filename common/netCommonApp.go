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
	"net/url"
	"time"
)

func ConnectionApp(
	startTimeOut time.Duration,
	stopTimeOut time.Duration,
	connectionName string,
	connectionInstancePrefix string,
	UseProxy bool,
	ProxyUrl *url.URL,
	ConnectionUrl *url.URL,
	stackName string,
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
		fx.Provide(func() *zap.Logger { return params.ZapLogger.Named(connectionName) }),
		fx.WithLogger(func(logger *zap.Logger) fxevent.Logger { return &fxevent.ZapLogger{Logger: logger} }),
		fx.Provide(fx.Annotated{Target: func() interfaces.IUniqueReferenceService { return params.UniqueSessionNumber_ }}),
		goCommsDefinitions.ProvideStringContext("ConnectionName", connectionName),
		goCommsDefinitions.ProvideStringContext("ConnectionInstancePrefix", connectionInstancePrefix),
		goCommsDefinitions.ProvideStringContext("StackName", stackName),
		goCommsDefinitions.ProvideUrl("ConnectionUrl", ConnectionUrl),
		goCommsDefinitions.ProvideUrl("ProxyUrl", ProxyUrl),
		goCommsDefinitions.ProvideBool("UseProxy", UseProxy),
		goConnectionManager.ProvideConnectionManager(params.ConnectionManager),
		goConnectionManager.ProvideCommandsToConnectionManager(),
		goConnectionManager.ProvideObtainConnectionManagerInformation(),
		goConnectionManager.ProvideRegisterToConnectionManager(),
		goConnectionManager.ProvidePublishConnectionInformation(),
		goCommsDefinitions.ProvideCancelContext(params.ParentContext),

		fx.Provide(fx.Annotated{Target: func() func() fx.Option {
			return additionalFxOptionsForConnectionInstance

		}}),

		fx.Options(option...),
	)
}
