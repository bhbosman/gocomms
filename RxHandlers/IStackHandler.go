package RxHandlers

import (
	"github.com/bhbosman/gocommon/model"
	"io"
)

type IStackHandler interface {
	io.Closer
	EmptyQueue()
	ClearCounters()
	PublishCounters(*model.PublishRxHandlerCounters)
}
