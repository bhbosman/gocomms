package common

import (
	"net"
	"sync"
)

type netConnWithSemaphoreWrapper struct {
	net.Conn
	mutex            sync.Mutex
	closed           bool
	releaseSemaphore interface {
		Release(int64)
	}
}

func NewNetConnWithSemaphoreWrapper(
	conn net.Conn,
	releaseSemaphore interface {
		Release(int64)
	}) *netConnWithSemaphoreWrapper {
	return &netConnWithSemaphoreWrapper{
		Conn:             conn,
		mutex:            sync.Mutex{},
		releaseSemaphore: releaseSemaphore,
	}
}

func (self *netConnWithSemaphoreWrapper) Close() error {
	if self.closed {
		return nil
	}
	self.mutex.Lock()
	defer self.mutex.Unlock()
	if self.closed {
		return nil
	}
	self.closed = true
	err := self.Conn.Close()
	local := self.releaseSemaphore
	self.releaseSemaphore = nil
	if local != nil {
		local.Release(1)
	}
	return err
}
