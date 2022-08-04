package RxHandlers

import (
	"context"
	"github.com/bhbosman/goprotoextra"
)

type FlatMapHandlerResult struct {
	UseDefaultPath bool
	Items          []interface{}
	RwsCount       int
	OtherCount     int
}

func NewFlatMapHandlerResult(
	useDefaultPath bool,
	items []interface{},
	RwsCount int,
	OtherCount int,
) FlatMapHandlerResult {
	return FlatMapHandlerResult{
		UseDefaultPath: useDefaultPath,
		Items:          items,
		OtherCount:     OtherCount,
		RwsCount:       RwsCount,
	}
}

type IRxMapStackHandler interface {
	IStackHandler
	MapReadWriterSize(context.Context, goprotoextra.ReadWriterSize) (goprotoextra.ReadWriterSize, error)
	ErrorState() error
	FlatMapHandler(item interface{}) (FlatMapHandlerResult, error)
}
