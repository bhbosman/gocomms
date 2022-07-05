package intf

import (
	"context"
	"github.com/bhbosman/gocommon/model"
	"go.uber.org/zap"
)

type IConnectionReactorFactory interface {
	Create(
		cancelCtx context.Context,
		cancelFunc context.CancelFunc,
		connectionCancelFunc model.ConnectionCancelFunc,
		logger *zap.Logger,
		userContext interface{}) (IConnectionReactor, error)
}
