package RxHandlers

import (
	"context"
	"github.com/bhbosman/gocommon/model"
	"github.com/reactivex/rxgo/v2"
	"go.uber.org/zap"
)

func All(
	description string,
	direction model.StreamDirection,
	Items chan rxgo.Item,
	logger *zap.Logger,
	ctx context.Context,
) (rxgo.NextFunc, rxgo.ErrFunc, rxgo.CompletedFunc, error) {
	b := false

	return func(data interface{}) {
			defer func() {
				if err := recover(); err != nil {
					logger.Error("panic in CreateSendData",
						zap.String("name", description),
						zap.Int("direction", int(direction)),
						zap.Any("recover", err),
					)
					b = false
				}
			}()
			if !b {
				item := rxgo.Of(data)
				item.SendContext(ctx, Items)
			}
		},
		func(err error) {
			defer func() {
				if err := recover(); err != nil {
					b = false
					logger.Error(
						"panic in CreateSendError",
						zap.String("name", description),
						zap.Int("direction", int(direction)),
						zap.Any("recover", err),
					)
				}
			}()
			if !b {
				item := rxgo.Error(err)
				item.SendContext(ctx, Items)
			}
		},
		func() {
			defer func() {
				err := recover()
				if err != nil {
					logger.Error(
						"panic in CreateComplete",
						zap.String("name", description),
						zap.Int("direction", int(direction)),
						zap.Any("recover", err),
					)
				}
			}()
			if !b {
				b = true
				close(Items)
			}
		}, nil
}
