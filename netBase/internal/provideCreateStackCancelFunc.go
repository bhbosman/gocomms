package internal

import (
	"github.com/bhbosman/gocommon/model"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"golang.org/x/net/context"
)

func ProvideCreateStackCancelFunc() fx.Option {
	return fx.Provide(
		func(
			params struct {
				fx.In
				CancelFunc context.CancelFunc
				Logger     *zap.Logger
			},
		) (model.ConnectionCancelFunc, error) {
			return func(context string, inbound bool, err error) {
				params.Logger.Error(context, zap.Error(err))
				params.CancelFunc()
			}, nil
		},
	)
}
