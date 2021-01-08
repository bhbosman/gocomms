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
	stackName string,
	connectionReactorFactories *ConnectionReactorFactories,
	ParentContext context.Context,
	StackFactory *TransportFactory,
	ConnectionManager connectionManager.IConnectionManager,
	connectionReactorFactoryName string,
	logFactory *gologging.Factory) fx.Option {
	return fx.Options(
		fx.Supply(connectionReactorFactories, StackFactory, logFactory),
		fx.Provide(fx.Annotated{Target: internal.CreateUrl(url)}),
		fx.Provide(fx.Annotated{Name: "StackName", Target: CreateStringContext(stackName)}),
		fx.Provide(fx.Annotated{Name: intf.ConnectionReactorFactoryName, Target: CreateStringContext(connectionReactorFactoryName)}),
		fx.Provide(fx.Annotated{Target: func() connectionManager.IConnectionManager { return ConnectionManager }}),
		fx.Provide(fx.Annotated{Target: func() connectionManager.IRegisterToConnectionManager { return ConnectionManager }}),
		fx.Provide(fx.Annotated{Target: func() connectionManager.IObtainConnectionManagerInformation { return ConnectionManager }}),
		fx.Provide(fx.Annotated{Target: func() connectionManager.ICommandsToConnectionManager { return ConnectionManager }}),
		fx.Provide(fx.Annotated{Target: func() rxgo.IPublishToConnectionManager { return ConnectionManager }}),
		fx.Provide(fx.Annotated{Target: func() (ctx context.Context, cancel context.CancelFunc) {
			return context.WithCancel(ParentContext)
		}}),
		fx.Provide(fx.Annotated{
			Target: func(
				params struct {
					fx.In
					StackName string `name:"StackName"`
					Factory   *TransportFactory
				}) (TransportFactoryFunction, error) {
				return params.Factory.Get(params.StackName)
			}}),
	)
}
