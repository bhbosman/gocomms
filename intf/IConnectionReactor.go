package intf

import (
	"github.com/reactivex/rxgo/v2"
	"io"
)

type IInitParams interface {
	OnSendToReactor() rxgo.NextFunc
	OnSendToConnection() rxgo.NextFunc
	NextFuncOutBoundChannel() rxgo.NextFunc
	NextFuncInBoundChannel() rxgo.NextFunc
}

type IConnectionReactor interface {
	io.Closer
	Init(params IInitParams) (rxgo.NextFunc, rxgo.ErrFunc, rxgo.CompletedFunc, error)
	Open() error
}

const ConnectionId = "ConnectionId"
