package RxHandlers

import (
	"context"
	"github.com/bhbosman/goprotoextra"
)

type IRxMapStackHandler interface {
	IStackHandler
	MapReadWriterSize(context.Context, goprotoextra.ReadWriterSize) (goprotoextra.ReadWriterSize, error)
	ErrorState() error
}
