package impl

import (
	"context"
	"fmt"
	"github.com/bhbosman/gocomms/connectionManager"
	"github.com/bhbosman/gocomms/internal"
	"github.com/bhbosman/gocomms/intf"
	"github.com/bhbosman/gologging"
	"github.com/bhbosman/goprotoextra"
	"github.com/reactivex/rxgo/v2"
	"go.uber.org/fx"
	"io"
	log2 "log"
	"math"
	"net"
	"net/url"

	"time"
)

type NetManager struct {
	connectionReactorFactories   *ConnectionReactorFactories
	CancelCtx                    context.Context
	cancelFunction               context.CancelFunc
	Logger                       *gologging.SubSystemLogger
	stackFactoryFunction         TransportFactoryFunction
	ConnectionReactorFactoryName string
	connectionManager            connectionManager.IConnectionManager
	Url                          *url.URL
	UserContext                  interface{}
	logFactory                   *gologging.Factory
	cfr                          intf.IConnectionReactorFactory
}

func (self *NetManager) NewConnectionInstance(connectionType internal.ConnectionType, conn net.Conn) (*fx.App, context.Context) {
	connectionName := fmt.Sprintf("Client Connection(%s-%s)", conn.RemoteAddr(), conn.LocalAddr())
	l := self.logFactory.Create(connectionName)
	var ctx context.Context
	fxApp := fx.New(
		fx.Populate(&ctx),
		fx.Logger(l),
		fx.StopTimeout(time.Hour),
		fx.Supply(l, self, self.connectionReactorFactories, self.stackFactoryFunction, self.Url, connectionType, self.logFactory),
		fx.Provide(fx.Annotated{Target: func() intf.IConnectionReactorFactory { return self.cfr }}),
		fx.Provide(fx.Annotated{Target: func() connectionManager.IConnectionManager { return self.connectionManager }}),
		fx.Provide(fx.Annotated{Target: func() connectionManager.IRegisterToConnectionManager { return self.connectionManager }}),
		fx.Provide(fx.Annotated{Target: func() connectionManager.IObtainConnectionManagerInformation { return self.connectionManager }}),
		fx.Provide(fx.Annotated{Target: func() connectionManager.ICommandsToConnectionManager { return self.connectionManager }}),
		fx.Provide(fx.Annotated{Target: func() rxgo.IPublishToConnectionManager { return self.connectionManager }}),
		fx.Provide(fx.Annotated{Target: provideNetConnection(conn)}),
		fx.Provide(fx.Annotated{Name: intf.ConnectionName, Target: CreateStringContext(connectionName)}),
		fx.Provide(fx.Annotated{Name: intf.UserContext, Target: func() interface{} { return self.UserContext }}),
		fx.Provide(fx.Annotated{Name: intf.ConnectionReactorFactoryName, Target: CreateStringContext(self.ConnectionReactorFactoryName)}),
		fx.Provide(fx.Annotated{Name: intf.ConnectionId, Target: CreateConnectionId()}),
		fx.Provide(func(netConnectionManager *NetManager) (context.Context, context.CancelFunc) {
			return context.WithCancel(netConnectionManager.CancelCtx)
		}),

		fx.Provide(fx.Annotated{Target: createClientContext}),
		fx.Provide(fx.Annotated{Target: createStackDefinition}),
		fx.Provide(fx.Annotated{Target: createStackCancelFunc}),
		fx.Provide(fx.Annotated{Target: createTransportLayer}),
		fx.Provide(fx.Annotated{Target: createToConnectionFunc}),
		fx.Provide(fx.Annotated{Target: createToReactorFunc}),
		fx.Provide(fx.Annotated{Target: createChannel}),

		fx.Invoke(invokeConnectionManager),
		fx.Invoke(invokeLogger),
		fx.Invoke(invokeChannel),
		fx.Invoke(invokeNetConnection),
		fx.Invoke(invokeObservable),
		fx.Invoke(invokeInboundTransportLayer),
		fx.Invoke(invokeOutBoundTransportLayer),
		fx.Invoke(invoke0007),
		fx.Invoke(invokeIndividualPipes),
		fx.Invoke(invokeIndividualStacks),
		fx.Invoke(invoke0008),
	)
	return fxApp, ctx
}

func NewNetManager(
	url *url.URL,
	connectionReactorFactories *ConnectionReactorFactories,
	cancelCtx context.Context,
	cancelFunction context.CancelFunc,
	logger *gologging.SubSystemLogger,
	stackFactoryFunction TransportFactoryFunction,
	ConnectionReactorFactoryName string,
	cfr intf.IConnectionReactorFactory,
	connectionManager connectionManager.IConnectionManager,
	logFactory *gologging.Factory,
	userContext interface{}) NetManager {
	return NetManager{
		connectionReactorFactories:   connectionReactorFactories,
		CancelCtx:                    cancelCtx,
		cancelFunction:               cancelFunction,
		Logger:                       logger,
		stackFactoryFunction:         stackFactoryFunction,
		ConnectionReactorFactoryName: ConnectionReactorFactoryName,
		connectionManager:            connectionManager,
		Url:                          url,
		UserContext:                  userContext,
		logFactory:                   logFactory,
		cfr:                          cfr,
	}
}

func createClientContext(
	params struct {
		fx.In
		Lifecycle                    fx.Lifecycle
		CancelCtx                    context.Context
		CancelFunc                   context.CancelFunc
		ConnectionManager            *ConnectionReactorFactories
		ConnectionName               string `name:"ConnectionName"`
		ConnectionReactorFactoryName string `name:"ConnectionReactorFactoryName"`
		Logger                       *gologging.SubSystemLogger
		Cfr                          intf.IConnectionReactorFactory
		ClientContext                interface{} `name:"UserContext"`
	}) (intf.IConnectionReactor, error) {
	params.Logger.LogWithLevel(0, func(logger *log2.Logger) {
		logger.Printf(fmt.Sprintf("createTransportLayer..."))
	})
	clientContext := params.Cfr.Create(
		params.ConnectionReactorFactoryName,
		params.CancelCtx,
		params.CancelFunc,
		params.Logger,
		params.ClientContext)
	var err error = nil
	if err != nil {
		params.Logger.LogWithLevel(0, func(logger *log2.Logger) {
			logger.Printf(fmt.Sprintf("createTransportLayer..."))
		})
		return nil, err
	}
	return clientContext, nil
}

func createStackCancelFunc(cancelFunc context.CancelFunc, logger *gologging.SubSystemLogger) internal.CancelFunc {
	return func(context string, inbound bool, err error) {
		_ = logger.ErrorWithDescription(context, err)
		cancelFunc()
	}
}

func createStackDefinition(
	params struct {
		fx.In
		ConnectionId         string `name:"ConnectionId"`
		Logger               *gologging.SubSystemLogger
		StackFactoryFunction TransportFactoryFunction
		CancelCtx            context.Context
		StackCancelFunc      internal.CancelFunc
		UserContext          intf.IConnectionReactor
		ConnectionManager    connectionManager.IConnectionManager
		ConnectionType       internal.ConnectionType
	}) (*internal.TwoWayPipeDefinition, error) {
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

func createChannel(
	params struct {
		fx.In
		CancelCtx    context.Context
		ConnectionId string `name:"ConnectionId"`
	}) (rxgo.Observable, *internal.ChannelManager) {
	result := make(chan rxgo.Item, 1024)
	obs := rxgo.FromChannel(result, rxgo.WithContext(params.CancelCtx))
	return obs, internal.NewChannelManager(result, "create channel", params.ConnectionId)

}

func invokeConnectionManager(
	params struct {
		fx.In
		ConnectionManager connectionManager.IConnectionManager
		CancelFunction    context.CancelFunc
		CancelCtx         context.Context
		ConnectionId      string `name:"ConnectionId"`
		Lifecycle         fx.Lifecycle
	}) {
	params.Lifecycle.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			if params.CancelCtx.Err() != nil {
				return params.CancelCtx.Err()
			}
			return params.ConnectionManager.RegisterConnection(
				params.ConnectionId,
				params.CancelFunction,
				params.CancelCtx)
		},
		OnStop: func(ctx context.Context) error {
			return params.ConnectionManager.DeregisterConnection(params.ConnectionId)
		},
	})
}

func invokeLogger(
	params struct {
		fx.In
		LifeCycle  fx.Lifecycle
		Logger     *gologging.SubSystemLogger
		CancelFunc context.CancelFunc
		CancelCtx  context.Context
	}) {

	params.LifeCycle.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			if params.CancelCtx.Err() != nil {
				return params.CancelCtx.Err()
			}
			return nil
		},
		OnStop: func(ctx context.Context) error {
			return nil
		},
	})
}

func invokeChannel(
	params struct {
		fx.In
		Channel    *internal.ChannelManager
		LifeCycle  fx.Lifecycle
		CancelFunc context.CancelFunc
		CancelCtx  context.Context
	}) {

	params.LifeCycle.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			return params.CancelCtx.Err()
		},
		OnStop: func(ctx context.Context) error {
			return params.Channel.Close()
		},
	})

}

func createToReactorFunc(
	params struct {
		fx.In
		CancelCtx context.Context
		Ch        *internal.ChannelManager
	}) goprotoextra.ToReactorFunc {
	return func(inline bool, any interface{}) error {
		err := params.CancelCtx.Err()
		if err != nil {
			return err
		}
		if !inline {
			go func(any interface{}) {
				params.Ch.Send(params.CancelCtx, rxgo.NewNextExternal(false, any))
			}(any)
			return nil
		}
		if params.Ch.Active() {
			return params.Ch.SendContextWithTimeOutAndRetries(
				rxgo.NewNextExternal(false, any),
				params.CancelCtx,
				time.Millisecond*500,
				5,
				params.Ch.Items)
		}
		return params.CancelCtx.Err()
	}
}
func createToConnectionFunc(
	params struct {
		fx.In
		TransportLayer *internal.TwoWayPipe
		CancelCtx      context.Context
	}) goprotoextra.ToConnectionFunc {
	return func(rw goprotoextra.ReadWriterSize) error {
		err := params.CancelCtx.Err()
		if err != nil {
			return err
		}
		return params.TransportLayer.SendOutgoingData(rw)
	}
}

func createTransportLayer(
	params struct {
		fx.In
		ConnectionId      string `name:"ConnectionId"`
		ConnectionManager connectionManager.IConnectionManager
		Lifecycle         fx.Lifecycle
		Logger            *gologging.SubSystemLogger
		CancelCtx         context.Context
		StackCancelFunc   internal.CancelFunc
		Def               *internal.TwoWayPipeDefinition
	}) (*internal.TwoWayPipe, error) {
	params.Logger.LogWithLevel(0, func(logger *log2.Logger) {
		logger.Printf(fmt.Sprintf("createTransportLayer..."))
	})
	transportLayer, err := params.Def.Build(
		params.ConnectionId,
		params.ConnectionManager,
		params.Logger,
		params.CancelCtx,
		params.StackCancelFunc)
	if err != nil {
		return nil, err
	}

	for _, pipeStart := range transportLayer.PipeState {
		localPipeStart := pipeStart
		params.Lifecycle.Append(fx.Hook{
			OnStart: func(ctx context.Context) error {
				if params.CancelCtx.Err() != nil {
					return params.CancelCtx.Err()
				}
				if localPipeStart.Start != nil {
					return localPipeStart.Start(params.CancelCtx)
				}
				return nil
			},
			OnStop: func(ctx context.Context) error {
				if localPipeStart.End != nil {
					return localPipeStart.End()
				}
				return nil
			},
		})
	}

	params.Lifecycle.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			return params.CancelCtx.Err()
		},
		OnStop: func(ctx context.Context) error {
			return transportLayer.Close()
		},
	})
	return transportLayer, nil
}

func provideNetConnection(conn net.Conn) func() (net.Conn, error) {
	return func() (net.Conn, error) {
		return conn, nil
	}
}

func invokeNetConnection(params struct {
	fx.In
	Conn       net.Conn
	Lifecycle  fx.Lifecycle
	CancelFunc context.CancelFunc
	CancelCtx  context.Context
}) {
	params.Lifecycle.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			return params.CancelCtx.Err()
		},
		OnStop: func(ctx context.Context) error {
			return params.Conn.Close()
		},
	})
}

func invokeObservable(
	params struct {
		fx.In
		Conn              net.Conn
		Url               *url.URL
		Lifecycle         fx.Lifecycle
		ClientContext     intf.IConnectionReactor
		ConnectionId      string `name:"ConnectionId"`
		ConnectionManager connectionManager.IConnectionManager
		ToConnectionFunc  goprotoextra.ToConnectionFunc
		ToReactorFunc     goprotoextra.ToReactorFunc
		Obs               rxgo.Observable
		CancelCtx         context.Context
	}) {
	params.Lifecycle.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			if params.CancelCtx.Err() != nil {
				return params.CancelCtx.Err()
			}
			onData, err := params.ClientContext.Init(
				params.Conn,
				params.Url,
				params.ConnectionId,
				params.ConnectionManager,
				params.ToConnectionFunc,
				params.ToReactorFunc)
			if err != nil {
				return err
			}
			params.Obs.(rxgo.InOutBoundObservable).DoNextExternal(
				-100,
				params.ConnectionId,
				"ConnectionReactor",
				rxgo.StreamDirectionInbound,
				params.ConnectionManager,
				onData)
			return nil
		},
		OnStop: func(ctx context.Context) error {
			return nil
		},
	})
}

func invokeInboundTransportLayer(
	params struct {
		fx.In
		Lifecycle         fx.Lifecycle
		TransportLayer    *internal.TwoWayPipe
		ChannelManager    *internal.ChannelManager
		CancelFunc        context.CancelFunc
		CancelCtx         context.Context
		ConnectionId      string `name:"ConnectionId"`
		ConnectionManager connectionManager.IConnectionManager
	}) {
	params.Lifecycle.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			if params.CancelCtx.Err() != nil {
				return params.CancelCtx.Err()
			}
			params.TransportLayer.InboundObservable.(rxgo.InOutBoundObservable).DoOnNextInOutBound(
				math.MinInt32,
				params.ConnectionId,
				"InboundTransportLayer",
				rxgo.StreamDirectionInbound,
				params.ConnectionManager,
				func(context context.Context, i goprotoextra.ReadWriterSize) {
					if context.Err() != nil {
						return
					}
					params.ChannelManager.Send(context, rxgo.NewNextExternal(true, i))
				})
			return nil
		},
		OnStop: func(ctx context.Context) error {
			return nil
		},
	})
}

func invokeOutBoundTransportLayer(
	params struct {
		fx.In
		Conn              net.Conn
		Lifecycle         fx.Lifecycle
		CancelFunc        context.CancelFunc
		TransportLayer    *internal.TwoWayPipe
		CancelCtx         context.Context
		ConnectionId      string `name:"ConnectionId"`
		ConnectionManager connectionManager.IConnectionManager
	}) {
	params.Lifecycle.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			if params.CancelCtx.Err() != nil {
				return params.CancelCtx.Err()
			}
			outboundNextClosed := false
			_ = params.TransportLayer.OutboundObservable.(rxgo.InOutBoundObservable).DoOnNextInOutBound(
				math.MaxInt16,
				params.ConnectionId,
				"Connection Write",
				rxgo.StreamDirectionOutbound,
				params.ConnectionManager,
				func(context context.Context, i goprotoextra.ReadWriterSize) {
					if context.Err() != nil {
						return
					}
					if outboundNextClosed {
						return
					}
					_, err := io.Copy(params.Conn, i)
					if err != nil {
						params.CancelFunc()
						return
					}
				})
			return nil
		},
		OnStop: func(ctx context.Context) error {
			return nil
		},
	})
}

func invokeIndividualStacks(
	params struct {
		fx.In
		Url                          *url.URL
		Lifecycle                    fx.Lifecycle
		TransportLayer               *internal.TwoWayPipe
		Conn                         net.Conn
		CancelCtx                    context.Context
		CancelFunc                   internal.CancelFunc
		ConnectionReactorFactories   *ConnectionReactorFactories
		ConnectionReactorFactoryName string `name:"ConnectionReactorFactoryName"`
		Cfr                          intf.IConnectionReactorFactory
	}) error {
	connectionReactorFactory := params.Cfr
	//connectionReactorFactory, ok := params.ConnectionReactorFactories.m[params.ConnectionReactorFactoryName]
	//if !ok {
	//	return fmt.Errorf("could not find ConnectionReactorFactory(%s) while invokeIndividualStacks", params.ConnectionReactorFactoryName)
	//}
	localConn := params.Conn
	for _, stackState := range params.TransportLayer.StackState {
		if params.CancelCtx.Err() != nil {
			return params.CancelCtx.Err()
		}
		localStackState := stackState
		params.Lifecycle.Append(fx.Hook{
			OnStart: func(ctx context.Context) error {
				if params.CancelCtx.Err() != nil {
					return params.CancelCtx.Err()
				}
				if localConn == nil {
					return nil
				}
				if localStackState.Start != nil {
					var err error
					localConn, err = localStackState.Start(
						internal.NewStackStartStateParams(
							localConn,
							params.Url,
							params.CancelCtx,
							params.CancelFunc,
							connectionReactorFactory))
					return err
				}
				return params.CancelCtx.Err()
			},
			OnStop: func(ctx context.Context) error {
				if localStackState.End != nil {
					var err error
					err = localStackState.End(internal.NewStackEndStateParams())
					return err
				}
				return nil
			},
		})
	}
	return nil
}

func invokeIndividualPipes(
	params struct {
		fx.In
		Lifecycle      fx.Lifecycle
		TransportLayer *internal.TwoWayPipe
		CancelFunc     context.CancelFunc
		CancelCtx      context.Context
	}) {
}
func invoke0007(
	params struct {
		fx.In
		Conn              net.Conn
		Lifecycle         fx.Lifecycle
		CancelFunc        internal.CancelFunc
		CancelCtx         context.Context
		TransportLayer    *internal.TwoWayPipe
		ConnectionManager connectionManager.IConnectionManager
		ConnectionId      string `name:"ConnectionId"`
	}) {
	params.Lifecycle.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			if params.CancelCtx.Err() != nil {
				return params.CancelCtx.Err()
			}
			go internal.ReadDataFromConnection(
				params.Conn,
				params.CancelFunc,
				params.CancelCtx,
				params.ConnectionManager,
				params.ConnectionId,
				1000000,
				"Connection Read",
				func(rws goprotoextra.IReadWriterSize, cancelCtx context.Context, CancelFunc internal.CancelFunc) {
					err := params.TransportLayer.ReceiveIncomingData(rws)
					if err != nil {
						CancelFunc("Connection Read", true, err)
						return
					}
				})
			return nil
		},
		OnStop: func(ctx context.Context) error {
			return nil
		},
	})
}

func invoke0008(
	params struct {
		fx.In
		Lifecycle     fx.Lifecycle
		ClientContext intf.IConnectionReactor
		CancelFunc    context.CancelFunc
		CancelCtx     context.Context
	}) {
	params.Lifecycle.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			if params.CancelCtx.Err() != nil {
				return params.CancelCtx.Err()
			}
			return params.ClientContext.Open()
		},
		OnStop: func(ctx context.Context) error {
			return params.ClientContext.Close()
		},
	})
}
