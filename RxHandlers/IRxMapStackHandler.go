package RxHandlers

import (
	"context"
	"github.com/bhbosman/goprotoextra"
	"github.com/reactivex/rxgo/v2"
)

type FlatMapHandlerResult struct {
	UseDefaultPath bool
	Items          rxgo.Observable
	MessageCount   int
}

func NewFlatMapHandlerResult(useDefaultPath bool, items rxgo.Observable, messageCount int) FlatMapHandlerResult {
	return FlatMapHandlerResult{UseDefaultPath: useDefaultPath, Items: items, MessageCount: messageCount}
}

type IRxMapStackHandler interface {
	IStackHandler
	MapReadWriterSize(context.Context, goprotoextra.ReadWriterSize) (goprotoextra.ReadWriterSize, error)
	ErrorState() error
	FlatMapHandler(item rxgo.Item) FlatMapHandlerResult
}
