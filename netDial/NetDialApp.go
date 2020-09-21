package netDial

import (
	"context"
	"fmt"
	"github.com/bhbosman/gocomms/connectionManager"
	"github.com/bhbosman/gocomms/impl"
	"github.com/bhbosman/gologging"
	"go.uber.org/fx"
)

type AppFuncInParams struct {
	fx.In
	ClientContextFactories *impl.ConnectionReactorFactories
	ParentContext          context.Context `name:"Application"`
	Lifecycle              fx.Lifecycle
	StackFactory           *impl.TransportFactory
	ConnectionManager      connectionManager.IConnectionManager
	LogFactory             *gologging.Factory
}
type AppFunc func(params AppFuncInParams) (*fx.App, error)

func NewNetDialApp(
	connectionName string,
	url string,
	stackName string,
	userContextFactoryName string,
	options ...DialAppSettingsApply) AppFunc {
	return func(params AppFuncInParams) (*fx.App, error) {
		return fx.New(
			fx.Supply(options),
			impl.CommonComponents(
				url,
				stackName,
				params.ClientContextFactories,
				params.ParentContext,
				params.StackFactory,
				params.ConnectionManager,
				userContextFactoryName,
				params.LogFactory),
			fx.Provide(
				func(params struct {
					fx.In
					Factory *gologging.Factory
				}) *gologging.SubSystemLogger {
					return params.Factory.Create(fmt.Sprintf("Dialer for %v", connectionName))
				}),
			fx.Provide(fx.Annotated{Target: newNetDialManager}),
			fx.Invoke(
				func(netManager *netDialManager, logger *gologging.SubSystemLogger, cancelFunction context.CancelFunc) {
					params.Lifecycle.Append(fx.Hook{
						OnStart: func(ctx context.Context) error {
							return netManager.Start(ctx)
						},
						OnStop: func(ctx context.Context) error {
							cancelFunction()
							return netManager.Stop(ctx)
						},
					})
				}),
		), nil
	}
}

