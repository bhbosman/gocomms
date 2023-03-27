package netBase

import (
	"context"
	"fmt"
	"github.com/bhbosman/goConn"
	"github.com/bhbosman/goConnectionManager"
	"github.com/bhbosman/gocommon/GoFunctionCounter"
	fx2 "github.com/bhbosman/gocommon/fx"
	"github.com/bhbosman/gocommon/messages"
	"github.com/bhbosman/gocommon/model"
	"github.com/bhbosman/gocommon/services/Providers"
	sss "github.com/bhbosman/gocommon/services/interfaces"
	"github.com/bhbosman/gocomms/intf"
	"github.com/bhbosman/gocomms/netBase/internal"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"io"
	"net"
	"net/url"
	"time"
)

func ProvideCancelContextWithRwc(cancelContext goConn.ICancellationContext) fx.Option {
	return fx.Provide(
		fx.Annotated{
			Target: func(
				params struct {
					fx.In
					ConnectionName          string `name:"ConnectionId"`
					Logger                  *zap.Logger
					PrimaryConnectionCloser io.Closer `name:"PrimaryConnection"`
				},
			) (context.Context, context.CancelFunc, goConn.ICancellationContext, error) {
				return cancelContext.CancelContext(), cancelContext.CancelFunc(), cancelContext, nil
			},
		},
	)
}

type ConnectionInstance struct {
	ConnectionUrl                            *url.URL
	UniqueSessionNumber                      sss.IUniqueReferenceService
	ConnectionManager                        goConnectionManager.IService
	CancelCtx                                goConn.ICancellationContext
	AdditionalFxOptionsForConnectionInstance func() fx.Option
	ZapLogger                                *zap.Logger
}

func NewConnectionInstance(
	connectionUrl *url.URL,
	uniqueSessionNumber sss.IUniqueReferenceService,
	connectionManager goConnectionManager.IService,
	cancelCtx goConn.ICancellationContext,
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
		ProvideCancelContextWithRwc(self.CancelCtx),
		self.AdditionalFxOptionsForConnectionInstance(),
		internal.ProvideRxOptions(),
		internal.ProvideCreateStackDefinition(),
		internal.ProvideCreateStackCancelFunc(),
		internal.ProvideChannel("InBoundChannel"),
		internal.ProvideChannel("OutBoundChannel"),
		internal.ProvideCreateToReactorFunc(),
		internal.ProvideCreateToConnectionFunc("ToConnectionFunc"),
		internal.ProvideCreateChannel(),
		internal.ProvideRxNextHandlerForNetConnRead22("net.conn.read"),
		internal.ProvideRxNextHandlerForNetConnWrite2("net.conn.write"),
		internal.ProvideConnReactorWrite2(),
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
		InvokeCancelContextWithRwc(),
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
) (messages.IApp, error) {
	fxAppOptions := self.NewConnectionInstanceOptions(
		uniqueReference,
		goFunctionCounter,
		connectionType,
		conn,
		newAddSettings(settingOptions...),
	)
	fxApp := fx.New(fxAppOptions)
	err := fxApp.Err()
	if err != nil {
		self.CancelCtx.CancelWithError("ddddd", err)
		//if conn != nil {
		//	// this is required, when providers fail, and the Conn.Close() invoker has not been added to the
		//	//fx.App lifecycle stack for destruction
		//	// this error can be ignored
		//	err = multierr.Append(err, conn.Close())
		//}
		self.ZapLogger.Error(
			"On a presumed error in the providers, the connection may double close",
			zap.Error(err))

		return nil, err
	}
	return fxApp, nil
}

// NewConnectionInstance is the default behaviour, which will use the StackName the connection was assigned with.
func (self ConnectionInstance) NewConnectionInstance(
	uniqueReference string,
	goFunctionCounter GoFunctionCounter.IService,
	connectionType model.ConnectionType,
	conn net.Conn,
) (messages.IApp, error) {
	return self.NewConnectionInstanceWithStackName(
		uniqueReference,
		goFunctionCounter,
		connectionType,
		conn,
	)
}
