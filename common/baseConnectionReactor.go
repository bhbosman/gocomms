package common

import (
	"context"
	"github.com/bhbosman/gocommon/model"
	"github.com/bhbosman/gomessageblock"
	"github.com/bhbosman/goprotoextra"
	"github.com/reactivex/rxgo/v2"
	"go.uber.org/zap"
)

type BaseConnectionReactor struct {
	// CancelCtx is the cancellation context associated with the connection. This can be used to check if the connection
	// have not been closed
	CancelCtx                      context.Context
	CancelFunc                     context.CancelFunc
	ConnectionCancelFunc           model.ConnectionCancelFunc
	Logger                         *zap.Logger
	ToConnection                   goprotoextra.ToConnectionFunc
	ToReactor                      goprotoextra.ToReactorFunc
	UserContext                    interface{}
	ToConnectionFuncReplacement    rxgo.NextFunc
	toConnectionReactorReplacement rxgo.NextFunc
}

func (self *BaseConnectionReactor) Init(
	toConnectionFunc goprotoextra.ToConnectionFunc,
	toConnectionReactor goprotoextra.ToReactorFunc,
	toConnectionFuncReplacement rxgo.NextFunc,
	toConnectionReactorReplacement rxgo.NextFunc,
) (rxgo.NextFunc, rxgo.ErrFunc, rxgo.CompletedFunc, error) {
	self.ToReactor = toConnectionReactor
	self.ToConnection = toConnectionFunc
	self.ToConnectionFuncReplacement = toConnectionFuncReplacement
	self.toConnectionReactorReplacement = toConnectionReactorReplacement
	return func(i interface{}) {

		}, func(err error) {

		}, func() {

		}, nil
}

func (self *BaseConnectionReactor) Close() error {
	return nil
}

func (self *BaseConnectionReactor) Open() error {
	return nil
}

func (self *BaseConnectionReactor) SendStringToConnection(s string) error {
	rws, err := gomessageblock.NewReaderWriterString(s)
	if err != nil {
		return err
	}
	return self.ToConnection(rws)
}

func NewBaseConnectionReactor(
	logger *zap.Logger,
	cancelCtx context.Context,
	cancelFunc context.CancelFunc,
	connectionCancelFunc model.ConnectionCancelFunc,
	userContext interface{}) BaseConnectionReactor {
	return BaseConnectionReactor{
		CancelCtx:            cancelCtx,
		CancelFunc:           cancelFunc,
		Logger:               logger,
		ToConnection:         nil,
		UserContext:          userContext,
		ConnectionCancelFunc: connectionCancelFunc,
	}
}
