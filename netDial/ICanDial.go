package netDial

import "sync/atomic"

type ICanDial interface {
	CanDial() bool
}

type CanDialDefaultImpl struct {
	count int64
}

func NewCanDialDefaultImpl() *CanDialDefaultImpl {
	return &CanDialDefaultImpl{}
}

func (self *CanDialDefaultImpl) CanDial() bool {
	return true
	//return self.count > 0
}

func (self *CanDialDefaultImpl) RemoveConsumer() {
	atomic.AddInt64(&self.count, -1)
}

func (self *CanDialDefaultImpl) AddConsumer() {
	atomic.AddInt64(&self.count, 1)
}
