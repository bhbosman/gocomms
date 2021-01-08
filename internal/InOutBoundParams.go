package internal

import "context"

type InOutBoundParams struct {
	Index   int
	Context context.Context
}

func NewInOutBoundParams(index int, context context.Context) InOutBoundParams {
	return InOutBoundParams{Index: index, Context: context}
}
