package RxHandlers

import (
	"github.com/bhbosman/goprotoextra"
	"io"
)

type IRxNextStackHandler interface {
	IStackHandler
	io.Closer
	OnError(err error)
	NextReadWriterSize(
		goprotoextra.ReadWriterSize,
		func(goprotoextra.ReadWriterSize) error,
		func(interface{}) error,
		func(int) error) error
	OnComplete()
	GetAdditionalBytesSend() int
	GetAdditionalBytesIncoming() int
}
