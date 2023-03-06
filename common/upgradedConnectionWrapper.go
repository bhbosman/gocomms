package common

import (
	"io"
	"net"
	"time"
)

type UpgradedConnectionWrapper struct {
	BytesRead  int
	BytesWrite int
	conn       net.Conn
}

func NewUpgradedConnectionWrapper(conn net.Conn) *UpgradedConnectionWrapper {
	return &UpgradedConnectionWrapper{
		conn: conn,
	}
}

func (self *UpgradedConnectionWrapper) Read(b []byte) (n int, err error) {
	if self.conn != nil {
		n, err = self.conn.Read(b)
		self.BytesRead += n
		return n, err
	}
	return 0, io.EOF
}

func (self *UpgradedConnectionWrapper) Write(b []byte) (n int, err error) {
	if self.conn != nil {
		n, err = self.conn.Write(b)
		return n, err
	}
	return 0, io.EOF
}

func (self *UpgradedConnectionWrapper) Close() error {
	if self.conn != nil {
		return self.conn.Close()
	}
	return nil
}

func (self *UpgradedConnectionWrapper) LocalAddr() net.Addr {
	return self.conn.LocalAddr()
}

func (self *UpgradedConnectionWrapper) RemoteAddr() net.Addr {
	return self.conn.RemoteAddr()
}

func (self *UpgradedConnectionWrapper) SetDeadline(t time.Time) error {
	return self.conn.SetDeadline(t)
}

func (self *UpgradedConnectionWrapper) SetReadDeadline(t time.Time) error {
	return self.conn.SetReadDeadline(t)
}

func (self *UpgradedConnectionWrapper) SetWriteDeadline(t time.Time) error {
	return self.conn.SetWriteDeadline(t)
}
