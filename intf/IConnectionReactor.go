package intf

import (
	"github.com/bhbosman/goprotoextra"
	"github.com/reactivex/rxgo/v2"
	"io"
)

type IConnectionReactor interface {
	io.Closer
	Init(
		onSend goprotoextra.ToConnectionFunc,
		toConnectionReactor goprotoextra.ToReactorFunc,
		onSendReplacement rxgo.NextFunc,
		toConnectionReactorReplacement rxgo.NextFunc,
	) (rxgo.NextFunc, rxgo.ErrFunc, rxgo.CompletedFunc, error)
	Open() error
}

const ConnectionId = "ConnectionId"
const UserContext = "UserContext"
