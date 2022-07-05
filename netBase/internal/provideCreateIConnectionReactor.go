package internal

import (
	"context"
	"github.com/bhbosman/gocommon/model"
	"github.com/bhbosman/gocomms/intf"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

// ProvideCreateIConnectionReactor deprecated
func ProvideCreateIConnectionReactor() fx.Option {
	return fx.Provide(
		fx.Annotated{
			Target: func(
				params struct {
					fx.In
					CancelCtx            context.Context
					CancelFunc           context.CancelFunc
					ConnectionCancelFunc model.ConnectionCancelFunc
					Logger               *zap.Logger
					Cfr                  intf.IConnectionReactorFactory
					ClientContext        interface{} `name:"UserContext"`
				},
			) (intf.IConnectionReactor, error) {
				params.Logger.Info("Creating Connection Reactor")
				return params.Cfr.Create(
					params.CancelCtx,
					params.CancelFunc,
					params.ConnectionCancelFunc,
					params.Logger,
					params.ClientContext)
			},
		},
	)
}
