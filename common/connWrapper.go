package common

import (
	"context"
	"github.com/bhbosman/goCommsDefinitions"
	"github.com/bhbosman/goerrors"
	"github.com/bhbosman/gomessageblock"
	"github.com/reactivex/rxgo/v2"
	"go.uber.org/multierr"
	"io"
	"net"
	"time"
)

type ConnWrapperNext func(payload []byte) (int, error)

type ConnWrapper struct {
	conn            goCommsDefinitions.ISpecificInformationForConnection
	isDisposed      bool
	cancelCtx       context.Context
	pipeReader      io.ReadCloser
	connWrapperNext ConnWrapperNext
	BytesWritten    int
	BytesRead       int
}

func (self *ConnWrapper) LocalAddr() net.Addr {
	return self.conn.LocalAddr()
}

func (self *ConnWrapper) RemoteAddr() net.Addr {
	return self.conn.RemoteAddr()
}

func (self *ConnWrapper) SetDeadline(t time.Time) error {
	return self.conn.SetDeadline(t)
}

func (self *ConnWrapper) SetReadDeadline(t time.Time) error {
	return self.conn.SetReadDeadline(t)
}

func (self *ConnWrapper) SetWriteDeadline(t time.Time) error {
	return self.conn.SetWriteDeadline(t)
}

func (self *ConnWrapper) Read(b []byte) (n int, err error) {
	if self.cancelCtx.Err() != nil {
		return 0, self.cancelCtx.Err()
	}
	n, err = self.pipeReader.Read(b)
	self.BytesRead += n
	return n, err
}

func (self *ConnWrapper) Write(b []byte) (n int, err error) {
	if self.cancelCtx.Err() != nil {
		return 0, self.cancelCtx.Err()
	}
	self.BytesWritten += len(b)
	return self.connWrapperNext(b)
}

func (self *ConnWrapper) Close() error {
	if !self.isDisposed {
		self.isDisposed = true
		var err error
		err = multierr.Append(err, self.pipeReader.Close())
		return err
	}
	return nil
}

func nextOutBoundPath(
	ctx context.Context,
	//itemChannel chan<- rxgo.Item,
	nextFunc rxgo.NextFunc,
) ConnWrapperNext {
	return func(b []byte) (n int, err error) {
		if ctx.Err() != nil {
			return 0, ctx.Err()
		}
		dataToConnection := gomessageblock.NewReaderWriterSize(len(b))
		if ctx.Err() != nil {
			return 0, ctx.Err()
		}
		n, err = dataToConnection.Write(b)
		if err != nil {
			return 0, err
		}
		if ctx.Err() != nil {
			return 0, ctx.Err()
		}
		//rxgo.Of(dataToConnection).SendContext(ctx, itemChannel)
		nextFunc(dataToConnection)
		return n, nil
	}
}

func NewConnWrapper(
	conn goCommsDefinitions.ISpecificInformationForConnection,
	cancelCtx context.Context,
	pipeRead io.ReadCloser,
	nextFunc rxgo.NextFunc,
	//itemChannel chan<- rxgo.Item,
) (*ConnWrapper, error) {
	var errList error = nil
	if conn == nil {
		errList = multierr.Append(errList, goerrors.InvalidParam)
	}
	if cancelCtx == nil {
		errList = multierr.Append(errList, goerrors.InvalidParam)
	}
	if pipeRead == nil {
		errList = multierr.Append(errList, goerrors.InvalidParam)
	}

	outBound := nextOutBoundPath(
		cancelCtx,
		//itemChannel,
		nextFunc,
	)
	if outBound == nil {
		errList = multierr.Append(errList, goerrors.InvalidParam)
	}

	if errList != nil {
		return nil, errList
	}

	return &ConnWrapper{
		conn:            conn,
		cancelCtx:       cancelCtx,
		pipeReader:      pipeRead,
		connWrapperNext: outBound,
	}, nil
}
