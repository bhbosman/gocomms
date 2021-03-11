package impl

import (
	"context"
	"github.com/bhbosman/gocomms/connectionManager"
	"github.com/bhbosman/gocomms/internal"
	"github.com/bhbosman/gocomms/intf"
	"github.com/bhbosman/gologging"
	"github.com/reactivex/rxgo/v2"
	"go.uber.org/fx"
)

func CommonComponents(
	url string,
	stackCreateFunction TransportFactoryFunction,
	connectionReactorFactories *ConnectionReactorFactories,
	ParentContext context.Context,
	ConnectionManager connectionManager.IConnectionManager,
	connectionReactorFactoryName string,
	cfr intf.IConnectionReactorFactory,
	logFactory *gologging.Factory) fx.Option {
	return fx.Options(
		fx.Supply(connectionReactorFactories, logFactory, stackCreateFunction),
		fx.Provide(fx.Annotated{Target: internal.CreateUrl(url)}),
		fx.Provide(fx.Annotated{Name: intf.ConnectionReactorFactoryName, Target: CreateStringContext(connectionReactorFactoryName)}),
		fx.Provide(fx.Annotated{Target: func() intf.IConnectionReactorFactory { return cfr }}),
		fx.Provide(fx.Annotated{Target: func() connectionManager.IConnectionManager { return ConnectionManager }}),
		fx.Provide(fx.Annotated{Target: func() connectionManager.IRegisterToConnectionManager { return ConnectionManager }}),
		fx.Provide(fx.Annotated{Target: func() connectionManager.IObtainConnectionManagerInformation { return ConnectionManager }}),
		fx.Provide(fx.Annotated{Target: func() connectionManager.ICommandsToConnectionManager { return ConnectionManager }}),
		fx.Provide(fx.Annotated{Target: func() rxgo.IPublishToConnectionManager { return ConnectionManager }}),
		fx.Provide(fx.Annotated{Target: func() (ctx context.Context, cancel context.CancelFunc) {
			return context.WithCancel(ParentContext)
		}}),
	)
}
