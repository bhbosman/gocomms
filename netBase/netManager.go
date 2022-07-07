package netBase

import (
	"context"
	"github.com/bhbosman/goCommsDefinitions"
	"github.com/bhbosman/goConnectionManager"
	"github.com/bhbosman/gocommon/GoFunctionCounter"
	"github.com/bhbosman/gocommon/Services/Providers"
	"github.com/bhbosman/gocommon/Services/interfaces"
	"github.com/bhbosman/gocommon/model"
	"github.com/bhbosman/gocomms/intf"
	"github.com/bhbosman/gocomms/netBase/internal"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"io"
	"net/url"
)

type NetManager struct {
	CancelCtx      context.Context
	cancelFunction context.CancelFunc
	ZapLogger      *zap.Logger
	//StackName                                string
	ConnectionManager                        goConnectionManager.IService
	ConnectionUrl                            *url.URL
	UseProxy                                 bool
	ProxyUrl                                 *url.URL
	UserContext                              interface{}
	Name                                     string
	ConnectionInstancePrefix                 string
	UniqueSessionNumber                      interfaces.IUniqueReferenceService
	AdditionalFxOptionsForConnectionInstance func() fx.Option
	GoFunctionCounter                        GoFunctionCounter.IService
}

func (self *NetManager) NewReaderWriterCloserInstanceOptions(
	uniqueReference string,
	goFunctionCounter GoFunctionCounter.IService,
	connectionType model.ConnectionType,
	rwc io.ReadWriteCloser,
	//stackName string,
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
		//goCommsDefinitions.ProvideStackName(stackName),
		fx.Supply(self.ConnectionUrl, connectionType),
		Providers.ProvideUniqueReferenceServiceInstance(self.UniqueSessionNumber),
		goConnectionManager.ProvideConnectionManager(self.ConnectionManager),
		goConnectionManager.ProvideCommandsToConnectionManager(),
		goConnectionManager.ProvideObtainConnectionManagerInformation(),
		goConnectionManager.ProvideRegisterToConnectionManager(),
		goConnectionManager.ProvidePublishConnectionInformation(),
		internal.ProvideReadWriteCloser(rwc),
		internal.ProvideExtractPipeInStates("PipeInStates"),
		internal.ProvideExtractPipeOutStates("PipeOutStates"),
		internal.ProvideCreateStackData(),
		fx.Provide(fx.Annotated{Name: intf.UserContext, Target: func() interface{} { return self.UserContext }}),
		fx.Provide(fx.Annotated{Name: intf.ConnectionId, Target: func() string { return uniqueReference }}),
		fx.Provide(fx.Annotated{Target: func() GoFunctionCounter.IService { return goFunctionCounter }}),
		goCommsDefinitions.ProvideCancelContext(self.CancelCtx),
		self.AdditionalFxOptionsForConnectionInstance(),
		internal.ProvideRxOptions(),
		internal.ProvideCreateStackDefinition(),
		internal.ProvideCreateStackCancelFunc(),
		internal.ProvideCreateTransportLayer00(),
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

func NewNetManager(
	name string,
	connectionInstancePrefix string,
	useProxy bool,
	proxyUrl *url.URL,
	connectionUrl *url.URL,
	cancelCtx context.Context,
	cancelFunction context.CancelFunc,
	//stackName string,
	connectionManager goConnectionManager.IService,
	userContext interface{},
	ZapLogger *zap.Logger,
	uniqueSessionNumber interfaces.IUniqueReferenceService,
	additionalFxOptionsForConnectionInstance func() fx.Option,
	GoFunctionCounter GoFunctionCounter.IService,
) (NetManager, error) {
	return NetManager{
		CancelCtx:      cancelCtx,
		cancelFunction: cancelFunction,
		ZapLogger:      ZapLogger,
		//StackName:                                stackName,
		ConnectionManager:                        connectionManager,
		ConnectionUrl:                            connectionUrl,
		ProxyUrl:                                 proxyUrl,
		UseProxy:                                 useProxy,
		UserContext:                              userContext,
		Name:                                     name,
		ConnectionInstancePrefix:                 connectionInstancePrefix,
		UniqueSessionNumber:                      uniqueSessionNumber,
		AdditionalFxOptionsForConnectionInstance: additionalFxOptionsForConnectionInstance,
		GoFunctionCounter:                        GoFunctionCounter,
	}, nil
}
