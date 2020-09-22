package internal

import (
	"context"
	"github.com/bhbosman/gorxextra"
	"github.com/reactivex/rxgo/v2"
	"sync"
	"time"
)

type ChannelManager struct {
	Items chan rxgo.Item
	Mutex *sync.Mutex
}

func (self *ChannelManager) Close() error {
	self.lock()
	defer self.unlock()
	close(self.Items)
	self.Items = nil
	return nil
}

func (self ChannelManager) lock() {
	self.Mutex.Lock()
}

func (self ChannelManager) unlock() {
	self.Mutex.Unlock()
}

func (self *ChannelManager) Active() bool {
	return self.Items != nil
}

func (self *ChannelManager) Send(ctx context.Context, i interface{}) {
	if self.Active() {
		item := rxgo.Of(i)
		self.lock()
		defer self.unlock()
		if self.Active() {
			item.SendContext(ctx, self.Items)
		}
	}
}

func (self *ChannelManager) SendError(ctx context.Context, i error) {
	if self.Active() {
		item := rxgo.Error(i)
		self.lock()
		defer self.unlock()
		if self.Active() {
			item.SendContext(ctx, self.Items)
		}
	}
}



func (self *ChannelManager) SendContextWithTimeOutAndRetries(
	i interface{},
	ctx context.Context,
	duration time.Duration,
	retryCount int,
	ch chan rxgo.Item) error{
	if self.Active() {
		item := rxgo.Of(i)
		self.lock()
		defer self.unlock()
		if self.Active() {
			return gorxextra.SendContextWithTimeOutAndRetries(
				item,
				ctx,
				time.Millisecond*500,
				5,
				self.Items)
		}
	}
	return ctx.Err()
}
