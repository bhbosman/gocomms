package internal

import (
	"github.com/bhbosman/gocommon/model"
	"github.com/bhbosman/gocomms/common"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

func ProvideCreateStackCancelFunc() fx.Option {
	return fx.Provide(
		func(
			params struct {
				fx.In
				CancellationContext common.ICancellationContext
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
