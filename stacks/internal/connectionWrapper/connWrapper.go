package connectionWrapper

import (
	"context"
	"io"
	"net"
)

type ConnWrapperNext func(payload []byte) (int, error)

type ConnWrapper struct {
	net.Conn
	cancelCtx       context.Context
	pipeReader      io.Reader
	connWrapperNext ConnWrapperNext
}

func (self *ConnWrapper) Read(b []byte) (n int, err error) {
	if self.cancelCtx.Err() != nil {
		return 0, self.cancelCtx.Err()
	}
	n, err = self.pipeReader.Read(b)
	return n, err
}

func (self *ConnWrapper) Write(b []byte) (n int, err error) {
	if self.cancelCtx.Err() != nil {
		return 0, self.cancelCtx.Err()
	}
	return self.connWrapperNext(b)
}

func (self ConnWrapper) Close() error {
	var err error = nil
	return err
}

func NewConnWrapper(
	conn net.Conn,
	cancelCtx context.Context,
	pipeRead io.Reader,
	outBound ConnWrapperNext) *ConnWrapper {
	return &ConnWrapper{
		Conn:            conn,
		cancelCtx:       cancelCtx,
		pipeReader:      pipeRead,
		connWrapperNext: outBound,
	}
}
