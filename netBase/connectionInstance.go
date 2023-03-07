package netBase

import (
	"fmt"
	"github.com/bhbosman/goCommsDefinitions"
	"github.com/bhbosman/goConnectionManager"
	"github.com/bhbosman/gocommon/GoFunctionCounter"
	sss "github.com/bhbosman/gocommon/Services/interfaces"
	fx2 "github.com/bhbosman/gocommon/fx"
	"github.com/bhbosman/gocommon/messages"
	"github.com/bhbosman/gocommon/model"
	"github.com/bhbosman/gocommon/services/Providers"
	"github.com/bhbosman/gocomms/intf"
	"github.com/bhbosman/gocomms/netBase/internal"
	"go.uber.org/fx"
	"go.uber.org/multierr"
	"go.uber.org/zap"
	"golang.org/x/net/context"
	"net"
	"net/url"
	"time"
)

type ConnectionInstance struct {
	ConnectionUrl                            *url.URL
	UniqueSessionNumber                      sss.IUniqueReferenceService
	ConnectionManager                        goConnectionManager.IService
	CancelCtx                                context.Context
	AdditionalFxOptionsForConnectionInstance func() fx.Option
	ZapLogger                                *zap.Logger
}

func NewConnectionInstance(
	connectionUrl *url.URL,
	uniqueSessionNumber sss.IUniqueReferenceService,
	connectionManager goConnectionManager.IService,
	cancelCtx context.Context,
	additionalFxOptionsForConnectionInstance func() fx.Option,
	zapLogger *zap.Logger,
) ConnectionInstance {
	return ConnectionInstance{
		ConnectionUrl:                            connectionUrl,
		UniqueSessionNumber:                      uniqueSessionNumber,
		ConnectionManager:                        connectionManager,
		CancelCtx:                                cancelCtx,
		AdditionalFxOptionsForConnectionInstance: additionalFxOptionsForConnectionInstance,
		ZapLogger:                                zapLogger,
	}
}

func (self ConnectionInstance) NewReaderWriterCloserInstanceOptions(
	uniqueReference string,
	goFunctionCounter GoFunctionCounter.IService,
	connectionType model.ConnectionType,
	conn net.Conn,
	provideLogger fx.Option,
	settingOptions ...INewConnectionInstanceSettingsApply,
) fx.Option {
	settings := &newConnectionInstanceSettings{
		options: nil,
	}
	for _, settingOption := range settingOptions {
		settingOption.apply(settings)
	}

	return fx.Options(
		internal.InvokeContextCancelFunc(),
		provideLogger,
		internal.WithLogger(),
		fx.Supply(self.ConnectionUrl, connectionType),
		Providers.ProvideUniqueReferenceServiceInstance(self.UniqueSessionNumber),
		goConnectionManager.ProvideConnectionManager(self.ConnectionManager),
		goConnectionManager.ProvideCommandsToConnectionManager(),
		goConnectionManager.ProvideObtainConnectionManagerInformation(),
		goConnectionManager.ProvideRegisterToConnectionManager(),
		goConnectionManager.ProvidePublishConnectionInformation(),
		internal.ProvideReadWriteCloser(conn),
		internal.ProvideExtractPipeInStates("PipeInStates"),
		internal.ProvideExtractPipeOutStates("PipeOutStates"),
		internal.ProvideCreateStackData(),
		fx.Provide(fx.Annotated{Name: intf.ConnectionId, Target: func() string { return uniqueReference }}),
		fx.Provide(fx.Annotated{Target: func() GoFunctionCounter.IService { return goFunctionCounter }}),
		goCommsDefinitions.ProvideCancelContextWithRwc(self.CancelCtx),
		self.AdditionalFxOptionsForConnectionInstance(),
		internal.ProvideRxOptions(),
		internal.ProvideCreateStackDefinition(),
		internal.ProvideCreateStackCancelFunc(),
		internal.ProvideChannel("InBoundChannel"),
		internal.ProvideCreateTransportLayer02Step01(),
		internal.ProvideChannel("OutBoundChannel"),
		internal.ProvideCreateTransportLayer02Step02(),
		internal.ProvideCreateToReactorFunc("ForReactor"),
		internal.ProvideCreateToConnectionFunc("ToConnectionFunc"),
		internal.ProvideCreateChannel("QWERTY"),
		internal.ProvideRxNextHandlerForNetConnRead22("net.conn.read"),
		internal.ProvideRxNextHandlerForNetConnWrite2("net.conn.write"),
		internal.ProvideConnReactorWrite2("conn.reactor.write"),
		internal.InvokeConnectionManager(),
		internal.InvokeCompleteIncomingObservable(),
		internal.InvokeInboundTransportLayer(),
		internal.InvokeOutBoundTransportLayer(),
		internal.InvokeFxLifeCycleReadDataFromConnectionStartStop(),
		internal.InvokeFxLifeCycleStackStateStartStop(),
		internal.InvokeFxLifeCyclePipeStateInStartStop(),
		internal.InvokeFxLifeCyclePipeStateOutStartStop(),
		internal.InvokeFxLifeCycleConnectionReactorStartStop(),
		internal.InvokeFxLifeCycleStartStopStackHeartbeat(),
		fx.Options(settings.options...),
		internal.InvokeContextCancelFunc(),
	)
}

func (self ConnectionInstance) NewConnectionInstanceOptions(
	uniqueReference string,
	goFunctionCounter GoFunctionCounter.IService,
	connectionType model.ConnectionType,
	conn net.Conn,
	settingOptions ...INewConnectionInstanceSettingsApply,
) fx.Option {
	rwcOptions := self.NewReaderWriterCloserInstanceOptions(
		uniqueReference,
		goFunctionCounter,
		connectionType,
		conn,
		fx.Provide(
			func() (*zap.Logger, error) {
				if self.ZapLogger == nil {
					return nil, fmt.Errorf("zap.Logger is nil. Please resolve")
				}
				return self.ZapLogger.With(
						zap.String("UniqueReference", uniqueReference),
						zap.String("LocalAddr", conn.LocalAddr().String()),
						zap.String("RemoteAddr", conn.RemoteAddr().String())),
					nil
			},
		),
		settingOptions...,
	)
	return fx2.NewFxApplicationOptions(
		time.Hour,
		time.Hour,
		rwcOptions,
		//fx.Provide(
		//	fx.Annotated{
		//		Target: func() (net.Conn, error) {
		//			if conn == nil {
		//				return nil, fmt.Errorf("connection is nil. Please resolve")
		//			}
		//			return conn, nil
		//		},
		//	},
		//),
	)
}

// NewConnectionInstanceWithStackName give the ability to override the stack name, after it has been set
// This is useful in SSH, when ssh can dynamically assign a stack name, depending on the ssh protocol required
func (self ConnectionInstance) NewConnectionInstanceWithStackName(
	uniqueReference string,
	goFunctionCounter GoFunctionCounter.IService,
	connectionType model.ConnectionType,
	conn net.Conn,
	settingOptions ...INewConnectionInstanceSettingsApply,
) (messages.IApp, context.Context, goCommsDefinitions.ICancellationContext, error) {
	var resultContext context.Context
	var resultCancelFunc context.CancelFunc
	var cancellationContext goCommsDefinitions.ICancellationContext
	fxAppOptions := self.NewConnectionInstanceOptions(
		uniqueReference,
		goFunctionCounter,
		connectionType,
		conn,
		NewAddFxOptions(fx.Populate(&resultContext)),
		NewAddFxOptions(fx.Populate(&resultCancelFunc, &cancellationContext)),
		newAddSettings(settingOptions...),
	)
	fxApp := fx.New(fxAppOptions)
	err := fxApp.Err()
	if err != nil {
		if resultCancelFunc != nil {
			resultCancelFunc()
		}
		if conn != nil {
			// this is required, when providers fail, and the Conn.Close() invoker has not been added to the
			//fx.App lifecycle stack for destruction
			// this error can be ignored
			err = multierr.Append(err, conn.Close())
		}
		if err != nil {
			self.ZapLogger.Error(
				"On a presumed error in the providers, the connection may double close",
				zap.Error(err))
		}
		return nil, nil, nil, err
	}
	return fxApp, resultContext, cancellationContext, nil
}

// NewConnectionInstance is the default behaviour, which will use the StackName the connection was assigned with.
func (self ConnectionInstance) NewConnectionInstance(
	uniqueReference string,
	goFunctionCounter GoFunctionCounter.IService,
	connectionType model.ConnectionType,
	conn net.Conn,
) (messages.IApp, context.Context, goCommsDefinitions.ICancellationContext, error) {
	return self.NewConnectionInstanceWithStackName(
		uniqueReference,
		goFunctionCounter,
		connectionType,
		conn,
	)
}
