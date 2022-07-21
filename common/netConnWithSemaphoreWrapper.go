package common

import (
	"github.com/bhbosman/goerrors"
	"net"
	"sync"
)

type netConnWithSemaphoreWrapper struct {
	net.Conn
	mutex           sync.Mutex
	closed          bool
	releaseCallback func()
}

func NewNetConnWithSemaphoreWrapper(
	conn net.Conn,
	releaseCallback func(),
) (*netConnWithSemaphoreWrapper, error) {
	if releaseCallback == nil {
		return nil, goerrors.InvalidParam
	}
	return &netConnWithSemaphoreWrapper{
		Conn:            conn,
		mutex:           sync.Mutex{},
		releaseCallback: releaseCallback,
	}, nil
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

	local := self.releaseCallback
	self.releaseCallback = nil
	if local != nil {
		local()
	}
	return err
}
