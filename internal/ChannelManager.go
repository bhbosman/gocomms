package internal

import (
	"context"
	"fmt"
	"github.com/bhbosman/gorxextra"
	"github.com/reactivex/rxgo/v2"
	"sync"
	"time"
)

type ChannelManager struct {
	Items chan rxgo.Item
	mutex *sync.Mutex
	Name  string
}

func NewChannelManager(items chan rxgo.Item, name, name2 string) *ChannelManager {
	return &ChannelManager{
		Items: items,
		mutex: &sync.Mutex{},
		Name:  fmt.Sprintf("%v %v", name, name2),
	}
}

func (self *ChannelManager) Close() error {
	self.lock()
	local := self.Items
	self.Items = nil
	self.unlock()
	if local != nil {
		close(local)
	}
	return nil
}

func (self ChannelManager) lock() {
	self.mutex.Lock()
}

func (self ChannelManager) unlock() {
	self.mutex.Unlock()
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
			n := len(self.Items)
			if n > 1000 {
				println("-->ChannelManager<--", self.Name, n)
			}
			if self.Active() {
				item.SendContext(ctx, self.Items)
			}
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
	ch chan rxgo.Item) error {
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
