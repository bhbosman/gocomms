package impl

import "net"

type netConnWithSemaphoreWrapper struct {
	net.Conn
	releaseSemaphore interface {
		Release(int64)
	}
}

func NewNetConnWithSemaphoreWrapper(
	conn net.Conn,
	releaseSemaphore interface {
		Release(int64)
	}) *netConnWithSemaphoreWrapper {
	return &netConnWithSemaphoreWrapper{Conn: conn, releaseSemaphore: releaseSemaphore}
}

func (self netConnWithSemaphoreWrapper) Close() error {
	err := self.Conn.Close()
	if self.releaseSemaphore != nil {
		self.releaseSemaphore.Release(1)
	}
	self.releaseSemaphore = nil
	return err
}
