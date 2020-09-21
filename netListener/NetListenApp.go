package netListener

import (
	"context"
	"fmt"

	"github.com/bhbosman/gocomms/connectionManager"
	"github.com/bhbosman/gocomms/impl"
	"github.com/bhbosman/gologging"
	"go.uber.org/fx"
	"net"
	url2 "net/url"
)

type NetListenAppFuncInParams struct {
	fx.In
	ClientContextFactories *impl.ConnectionReactorFactories
	ParentContext          context.Context `name:"Application"`
	Lifecycle              fx.Lifecycle
	StackFactory           *impl.TransportFactory
	ConnectionManager      connectionManager.IConnectionManager
	LogFactory             *gologging.Factory
}

func NewNetListenApp(
	connectionName string,
	url string,
	stackName string,
	userContextFactoryName string, settings ...ListenAppSettingsApply) NewNetListenAppFunc {
	return func(params NetListenAppFuncInParams) (*fx.App, error) {
		return fx.New(
			fx.Supply(settings),
			impl.CommonComponents(
				url,
				stackName,
				params.ClientContextFactories,
				params.ParentContext,
				params.StackFactory,
				params.ConnectionManager,
				userContextFactoryName,
				params.LogFactory),

			fx.Provide(fx.Annotated{Target: newNetListenManager}),
			fx.Provide(
				func(Lifecycle fx.Lifecycle, url *url2.URL) (net.Listener, error) {
					con, err := net.Listen(url.Scheme, url.Host)
					if err != nil {
						return nil, err
					}
					Lifecycle.Append(fx.Hook{
						OnStart: nil,
						OnStop: func(ctx context.Context) error {
							return con.Close()
						},
					})
					return con, nil
				}),
			fx.Provide(
				func(params struct {
					fx.In
					Factory *gologging.Factory
				}) *gologging.SubSystemLogger {
					return params.Factory.Create(fmt.Sprintf("Listener for %v", connectionName))
				}),

			fx.Invoke(
				func(netManager *netListenManager, logger *gologging.SubSystemLogger, cancelFunc context.CancelFunc) {
					params.Lifecycle.Append(fx.Hook{
						OnStart: func(ctx context.Context) error {
							netManager.listenForNewConnections()
							return nil
						},
						OnStop: func(ctx context.Context) error {
							cancelFunc()
							return nil
						},
					})
				}),
		), nil
	}
}
