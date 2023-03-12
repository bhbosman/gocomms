package RxHandlers

import (
	"context"
	"github.com/bhbosman/goCommsDefinitions"
	"github.com/bhbosman/gocommon/model"
	"github.com/reactivex/rxgo/v2"
	"go.uber.org/zap"
)

func All2(
	description string,
	direction model.StreamDirection,
	items chan<- rxgo.Item,
	logger *zap.Logger,
	ctx context.Context,
	useCompleteCallback bool,
) (goCommsDefinitions.IRxNextHandler, error) {
	i, t, e, c, active, err := all(description, direction, items, logger, ctx)

	if !useCompleteCallback {
		c = nil
	}
	if err != nil {
		return nil, err
	}
	return goCommsDefinitions.NewDefaultRxNextHandler(i, t, e, c, active)
}

func all(
	description string,
	direction model.StreamDirection,
	items chan<- rxgo.Item,
	logger *zap.Logger,
	ctx context.Context,
) (rxgo.NextFunc, goCommsDefinitions.TryNextFunc, rxgo.ErrFunc, rxgo.CompletedFunc, goCommsDefinitions.IsNextActive, error) {
	isClosed := false
	//errorHappen := false
	return func(data interface{}) {
			defer func() {
				if err := recover(); err != nil {
					logger.Error("panic in CreateSendData",
						zap.String("name", description),
						zap.Int("direction", int(direction)),
						zap.Any("recover", err),
					)
					//errorHappen = true
				}
			}()
			if !isClosed {
				item := rxgo.Of(data)
				item.SendContext(ctx, items)
			}
		},
		func(data interface{}) bool {
			defer func() {
				if err := recover(); err != nil {
					logger.Error("panic in CreateSendData",
						zap.String("name", description),
						zap.Int("direction", int(direction)),
						zap.Any("recover", err),
					)
					//errorHappen = true
				}
			}()
			if !isClosed {
				item := rxgo.Of(data)
				return item.SendNonBlocking(items)
			}
			return false
		},
		func(err error) {
			defer func() {
				if err := recover(); err != nil {
					logger.Error(
						"panic in CreateSendError",
						zap.String("name", description),
						zap.Int("direction", int(direction)),
						zap.Any("recover", err),
					)
					//errorHappen = true
				}
			}()
			if !isClosed {
				item := rxgo.Error(err)
				item.SendContext(ctx, items)
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
					//errorHappen = true
				}
			}()
			if !isClosed {
				isClosed = true
				close(items)
			}
		},
		func() bool {
			return !isClosed
		}, nil
}
