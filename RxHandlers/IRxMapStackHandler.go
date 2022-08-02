package RxHandlers

import (
	"context"
	"github.com/bhbosman/goprotoextra"
)

type FlatMapHandlerResult struct {
	UseDefaultPath bool
	Items          []interface{}
	MessageCount   int
}

func NewFlatMapHandlerResult(useDefaultPath bool, items []interface{}, messageCount int) FlatMapHandlerResult {
	return FlatMapHandlerResult{UseDefaultPath: useDefaultPath, Items: items, MessageCount: messageCount}
}

type IRxMapStackHandler interface {
	IStackHandler
	MapReadWriterSize(context.Context, goprotoextra.ReadWriterSize) (goprotoextra.ReadWriterSize, error)
	ErrorState() error
	FlatMapHandler(item interface{}) FlatMapHandlerResult
}
