package common

import (
	"context"
	"github.com/bhbosman/gocommon/GoFunctionCounter"
	"github.com/bhbosman/gocommon/model"
	"github.com/bhbosman/gocommon/pubSub"
	"github.com/cskr/pubsub"
	"github.com/reactivex/rxgo/v2"
	"go.uber.org/multierr"
	"go.uber.org/zap"
)

type BaseConnectionReactor struct {
	// CancelCtx is the cancellation context associated with the connection. This can be used to check if the connection
	// have not been closed
	CancelCtx                   context.Context
	CancelFunc                  context.CancelFunc
	ConnectionCancelFunc        model.ConnectionCancelFunc
	Logger                      *zap.Logger
	OnSendToReactor             rxgo.NextFunc
	OnSendToConnection          rxgo.NextFunc
	UniqueReference             string
	PubSub                      *pubsub.PubSub
	GoFunctionCounter           GoFunctionCounter.IService
	OnSendToReactorPubSubBag    *pubsub.NextFuncSubscription
	OnSendToConnectionPubSubBag *pubsub.NextFuncSubscription
}

func (self *BaseConnectionReactor) Init(
	onSendToReactor rxgo.NextFunc,
	onSendToConnection rxgo.NextFunc,
) (rxgo.NextFunc, rxgo.ErrFunc, rxgo.CompletedFunc, error) {
	self.OnSendToReactor = onSendToReactor
	self.OnSendToConnection = onSendToConnection
	self.OnSendToReactorPubSubBag = pubsub.NewNextFuncSubscription(self.OnSendToReactor)
	self.OnSendToConnectionPubSubBag = pubsub.NewNextFuncSubscription(self.OnSendToConnection)

	return func(i interface{}) {

		}, func(err error) {

		}, func() {

		},
		nil
}

func (self *BaseConnectionReactor) Close() error {
	err := pubSub.Unsubscribe(
		"Unsubscribe from Pubsub",
		self.PubSub,
		self.GoFunctionCounter,
		self.OnSendToReactorPubSubBag,
	)
	err = multierr.Append(
		err,
		pubSub.Unsubscribe(
			"Unsubscribe from Pubsub",
			self.PubSub,
			self.GoFunctionCounter,
			self.OnSendToConnectionPubSubBag,
		),
	)
	return err
}

func (self *BaseConnectionReactor) Open() error {
	return nil
}

func NewBaseConnectionReactor(
	logger *zap.Logger,
	cancelCtx context.Context,
	cancelFunc context.CancelFunc,
	connectionCancelFunc model.ConnectionCancelFunc,
	uniqueReference string,
	PubSub *pubsub.PubSub,
	GoFunctionCounter GoFunctionCounter.IService,
) BaseConnectionReactor {
	return BaseConnectionReactor{
		CancelCtx:            cancelCtx,
		CancelFunc:           cancelFunc,
		ConnectionCancelFunc: connectionCancelFunc,
		Logger:               logger,
		UniqueReference:      uniqueReference,
		PubSub:               PubSub,
		GoFunctionCounter:    GoFunctionCounter,
	}
}
