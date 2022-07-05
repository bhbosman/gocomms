package internal

import (
	"fmt"
	"go.uber.org/fx"
	"go.uber.org/fx/fxevent"
	"go.uber.org/zap"
)

func WithLogger() fx.Option {
	return fx.WithLogger(
		func(
			params struct {
				fx.In
				Logger *zap.Logger
			},
		) (fxevent.Logger, error) {
			if params.Logger == nil {
				return nil, fmt.Errorf("logger *zap.Logger is nil. Please resolve")
			}
			return &fxevent.ZapLogger{
				Logger: params.Logger,
			}, nil
		},
	)
}
