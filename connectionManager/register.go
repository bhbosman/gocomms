package connectionManager

import (
	"context"
	"github.com/cskr/pubsub"
	"github.com/reactivex/rxgo/v2"
	"go.uber.org/fx"
)

func RegisterDefaultConnectionManager() fx.Option {
	return fx.Options(
		fx.Provide(
			func(params struct {
				fx.In
				CancelCtx context.Context `name:"Application"`
				PubSub    *pubsub.PubSub  `name:"Application"`
			}) (IConnectionManager, IRegisterToConnectionManager, IObtainConnectionManagerInformation, ICommandsToConnectionManager, rxgo.IPublishToConnectionManager) {
				cm := NewConnectionManager(params.PubSub, params.CancelCtx)
				return cm, cm, cm, cm, cm
			}),
		fx.Invoke(func(params struct {
			fx.In
			LifeCycle         fx.Lifecycle
			ConnectionManager IConnectionManager
		}) {
			params.LifeCycle.Append(fx.Hook{
				OnStart: params.ConnectionManager.Start,
				OnStop:  params.ConnectionManager.Stop,
			})
		}),
	)
}
