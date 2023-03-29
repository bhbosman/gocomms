package internal

import (
	"github.com/bhbosman/gocommon"
	"github.com/bhbosman/gocommon/model"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

func ProvideCreateStackCancelFunc() fx.Option {
	return fx.Provide(
		func(
			params struct {
				fx.In
				CancellationContext gocommon.ICancellationContext
				Logger              *zap.Logger
			},
		) (model.ConnectionCancelFunc, error) {
			return func(context string, inbound bool, err error) {
				params.Logger.Error(context, zap.Error(err))
				params.CancellationContext.CancelWithError("ProvideCreateStackCancelFunc", err)
			}, nil
		},
	)
}
