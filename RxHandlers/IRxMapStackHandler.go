package RxHandlers

import (
	"context"
)

type FlatMapHandlerResult struct {
	UseDefaultPath bool
	Items          []interface{}
	RwsCount       int
	OtherCount     int
	BytesIn        int
	BytesOut       int
}

func NewFlatMapHandlerResult(
	useDefaultPath bool,
	items []interface{},
	RwsCount int,
	OtherCount int,
	BytesIn int,
	BytesOut int,
) FlatMapHandlerResult {
	return FlatMapHandlerResult{
		UseDefaultPath: useDefaultPath,
		Items:          items,
		RwsCount:       RwsCount,
		OtherCount:     OtherCount,
		BytesIn:        BytesIn,
		BytesOut:       BytesOut,
	}
}

type IRxMapStackHandler interface {
	IStackHandler
	MapReadWriterSize(context.Context, interface{}) (interface{}, error)
	FlatMapHandler(ctx context.Context, item interface{}) (FlatMapHandlerResult, error)
}
