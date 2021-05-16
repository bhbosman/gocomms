package impl

import (
	"context"
	"fmt"
	"github.com/bhbosman/gocomms/connectionManager"
	internal2 "github.com/bhbosman/gocomms/impl/internal"
	commsInternal "github.com/bhbosman/gocomms/internal"
	"github.com/bhbosman/gocomms/intf"
	"github.com/bhbosman/gologging"
	"github.com/reactivex/rxgo/v2"
	"go.uber.org/fx"
	log2 "log"
	"net"
	"net/url"

	"time"
)

type NetManager struct {
	CancelCtx            context.Context
	cancelFunction       context.CancelFunc
	Logger               *gologging.SubSystemLogger
	stackFactoryFunction TransportFactoryFunction
	connectionManager    connectionManager.IConnectionManager
	Url                  *url.URL
	UserContext          interface{}
	logFactory           *gologging.Factory
	cfr                  intf.IConnectionReactorFactory
}

func (self *NetManager) NewConnectionInstance(connectionType commsInternal.ConnectionType, conn net.Conn) (*fx.App, context.Context) {

	connectionName := fmt.Sprintf("Client Connection(%s-%s)", conn.RemoteAddr(), conn.LocalAddr())
	l := self.logFactory.Create(connectionName)
	var ctx context.Context
	fxApp := fx.New(
		fx.Populate(&ctx),
		fx.Logger(l),
		fx.StopTimeout(time.Hour),
		fx.Supply(l, self, self.stackFactoryFunction, self.Url, connectionType, self.logFactory),
		fx.Provide(fx.Annotated{Target: func() intf.IConnectionReactorFactory { return self.cfr }}),
		fx.Provide(fx.Annotated{Target: func() connectionManager.IConnectionManager { return self.connectionManager }}),
		fx.Provide(fx.Annotated{Target: func() connectionManager.IRegisterToConnectionManager { return self.connectionManager }}),
		fx.Provide(fx.Annotated{Target: func() connectionManager.IObtainConnectionManagerInformation { return self.connectionManager }}),
		fx.Provide(fx.Annotated{Target: func() connectionManager.ICommandsToConnectionManager { return self.connectionManager }}),
		fx.Provide(fx.Annotated{Target: func() rxgo.IPublishToConnectionManager { return self.connectionManager }}),
		fx.Provide(fx.Annotated{Target: internal2.ProvideNetConnection(conn)}),
		fx.Provide(fx.Annotated{Name: "PipeInData", Target: internal2.CreatePipeInStateData}),
		fx.Provide(fx.Annotated{Name: "PipeOutData", Target: internal2.CreatePipeOutStateData}),
		fx.Provide(fx.Annotated{Name: "PipeInStates", Target: extractPipeInStates}),
		fx.Provide(fx.Annotated{Name: "PipeOutStates", Target: extractPipeOutStates}),
		fx.Provide(fx.Annotated{Name: "StackData", Target: internal2.CreateStackData}),
		fx.Provide(fx.Annotated{Name: intf.ConnectionName, Target: CreateStringContext(connectionName)}),
		fx.Provide(fx.Annotated{Name: intf.UserContext, Target: func() interface{} { return self.UserContext }}),
		fx.Provide(fx.Annotated{Name: intf.ConnectionId, Target: CreateConnectionId()}),
		fx.Provide(func(netConnectionManager *NetManager) (context.Context, context.CancelFunc) {
			return context.WithCancel(netConnectionManager.CancelCtx)
		}),

		fx.Provide(fx.Annotated{Target: internal2.CreateClientContext}),
		fx.Provide(fx.Annotated{Target: createStackDefinition}),
		fx.Provide(fx.Annotated{Target: internal2.CreateStackCancelFunc}),
		fx.Provide(fx.Annotated{Target: internal2.CreateTransportLayer00}),
		fx.Provide(fx.Annotated{Target: internal2.CreateTransportLayer02Step01}),
		fx.Provide(fx.Annotated{Target: internal2.CreateTransportLayer02Step02}),
		fx.Provide(fx.Annotated{Target: internal2.CreateToConnectionFunc}),
		fx.Provide(fx.Annotated{Target: internal2.CreateToReactorFunc}),
		fx.Provide(fx.Annotated{Target: internal2.CreateChannel}),

		fx.Invoke(internal2.InvokeConnectionManager),
		fx.Invoke(internal2.InvokeLogger),
		fx.Invoke(internal2.InvokeChannel),
		fx.Invoke(internal2.InvokeNetConnection),

		fx.Invoke(internal2.InvokeCompleteIncomingObservable),

		fx.Invoke(internal2.InvokeInboundTransportLayer),
		fx.Invoke(internal2.InvokeOutBoundTransportLayer),

		fx.Invoke(internal2.Invoke0007),
		fx.Invoke(internal2.InvokeTwoWayPipe),
		fx.Invoke(internal2.Invoke0008),
		fx.Invoke(internal2.InvokeStackStateStartStop),
		fx.Invoke(internal2.InvokePipeStateInStartStop),
		fx.Invoke(internal2.InvokePipeStateOutStartStop),
	)
	return fxApp, ctx
}

func NewNetManager(
	url *url.URL,
	cancelCtx context.Context,
	cancelFunction context.CancelFunc,
	logger *gologging.SubSystemLogger,
	stackFactoryFunction TransportFactoryFunction,
	cfr intf.IConnectionReactorFactory,
	connectionManager connectionManager.IConnectionManager,
	logFactory *gologging.Factory,
	userContext interface{}) NetManager {
	return NetManager{
		CancelCtx:            cancelCtx,
		cancelFunction:       cancelFunction,
		Logger:               logger,
		stackFactoryFunction: stackFactoryFunction,
		connectionManager:    connectionManager,
		Url:                  url,
		UserContext:          userContext,
		logFactory:           logFactory,
		cfr:                  cfr,
	}
}

func createStackDefinition(
	params struct {
		fx.In
		ConnectionId         string `name:"ConnectionId"`
		Logger               *gologging.SubSystemLogger
		StackFactoryFunction TransportFactoryFunction
		CancelCtx            context.Context
		StackCancelFunc      commsInternal.CancelFunc
		UserContext          intf.IConnectionReactor
		ConnectionManager    connectionManager.IConnectionManager
		ConnectionType       commsInternal.ConnectionType
	}) (*commsInternal.TwoWayPipeDefinition, error) {
	params.Logger.LogWithLevel(0, func(logger *log2.Logger) {
		logger.Printf(fmt.Sprintf("createStackDefinition..."))
	})

	function, err := params.StackFactoryFunction(
		params.ConnectionType,
		params.ConnectionId,
		params.UserContext,
		params.ConnectionManager,
		params.CancelCtx,
		params.StackCancelFunc,
		rxgo.WithContext(params.CancelCtx),
		rxgo.WithBufferedChannel(1024),
		rxgo.WithBackPressureStrategy(rxgo.Block))
	return function, err
}

func extractPipeInStates(params struct {
	fx.In
	LifeCycle            fx.Lifecycle
	TwoWayPipeDefinition *commsInternal.TwoWayPipeDefinition
}) ([]*commsInternal.PipeState, error) {
	return params.TwoWayPipeDefinition.BuildInBoundPipeStates()
}

func extractPipeOutStates(params struct {
	fx.In
	LifeCycle            fx.Lifecycle
	TwoWayPipeDefinition *commsInternal.TwoWayPipeDefinition
}) ([]*commsInternal.PipeState, error) {
	return params.TwoWayPipeDefinition.BuildOutBoundPipeStates()
}
