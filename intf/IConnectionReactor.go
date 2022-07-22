package intf

import (
	"github.com/reactivex/rxgo/v2"
	"io"
)

type IConnectionReactor interface {
	io.Closer
	Init(
		onSendToReactor rxgo.NextFunc,
		onSendToConnection rxgo.NextFunc,
	) (rxgo.NextFunc, rxgo.ErrFunc, rxgo.CompletedFunc, error)
	Open() error
}

const ConnectionId = "ConnectionId"

//const UserContext = "UserContext"
