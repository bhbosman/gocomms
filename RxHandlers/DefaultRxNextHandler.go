package RxHandlers

import "github.com/reactivex/rxgo/v2"

type DefaultRxNextHandler struct {
	i rxgo.NextFunc
	e rxgo.ErrFunc
	c rxgo.CompletedFunc
}

func NewDefaultRxNextHandler(i rxgo.NextFunc, e rxgo.ErrFunc, c rxgo.CompletedFunc) *DefaultRxNextHandler {
	return &DefaultRxNextHandler{
		i: i,
		e: e,
		c: c,
	}
}

func (self *DefaultRxNextHandler) OnSendData(i interface{}) {
	self.i(i)
}

func (self *DefaultRxNextHandler) OnError(err error) {
	self.e(err)
}

func (self *DefaultRxNextHandler) OnComplete() {
	self.c()
}
