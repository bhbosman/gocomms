package RxHandlers

import (
	"github.com/bhbosman/gocommon/model"
)

type IStackHandler interface {
	EmptyQueue()
	ClearCounters()
	PublishCounters(*model.PublishRxHandlerCounters)
}
