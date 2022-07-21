package common

import (
	"context"
	"github.com/bhbosman/gocommon/model"
	"github.com/reactivex/rxgo/v2"
	"go.uber.org/zap"
)

type BaseConnectionReactor struct {
	// CancelCtx is the cancellation context associated with the connection. This can be used to check if the connection
	// have not been closed
	CancelCtx            context.Context
	CancelFunc           context.CancelFunc
	ConnectionCancelFunc model.ConnectionCancelFunc
	Logger               *zap.Logger
	OnSendToReactor      rxgo.NextFunc
	OnSendToConnection   rxgo.NextFunc
	UniqueReference      string
}

func (self *BaseConnectionReactor) Init(
	onSendToReactor rxgo.NextFunc,
	onSendToConnection rxgo.NextFunc,
) (rxgo.NextFunc, rxgo.ErrFunc, rxgo.CompletedFunc, chan interface{}, error) {
	self.OnSendToReactor = onSendToReactor
	self.OnSendToConnection = onSendToConnection
	return func(i interface{}) {

		}, func(err error) {

		}, func() {

		},
		nil,
		nil
}

func (self *BaseConnectionReactor) Close() error {
	return nil
}

func (self *BaseConnectionReactor) Open() error {
	return nil
}

//
//func (self *BaseConnectionReactor) SendStringToConnection(s string) error {
//	rws, err := gomessageblock.NewReaderWriterString(s)
//	if err != nil {
//		return err
//	}
//	return self.ToConnection(rws)
//}

func NewBaseConnectionReactor(
	logger *zap.Logger,
	cancelCtx context.Context,
	cancelFunc context.CancelFunc,
	connectionCancelFunc model.ConnectionCancelFunc,
	uniqueReference string,
) BaseConnectionReactor {
	return BaseConnectionReactor{
		CancelCtx:            cancelCtx,
		CancelFunc:           cancelFunc,
		Logger:               logger,
		ConnectionCancelFunc: connectionCancelFunc,
		UniqueReference:      uniqueReference,
	}
}
