package intf

import (
	"github.com/reactivex/rxgo/v2"
	"io"
)

type IInitParams interface {
	OnSendToReactor() rxgo.NextFunc
	OnSendToConnection() rxgo.NextFunc
}

type initParams struct {
	onSendToReactor    rxgo.NextFunc
	onSendToConnection rxgo.NextFunc
}

func NewInitParams(onSendToReactor rxgo.NextFunc, onSendToConnection rxgo.NextFunc) IInitParams {
	return &initParams{onSendToReactor: onSendToReactor, onSendToConnection: onSendToConnection}
}

func (self *initParams) OnSendToReactor() rxgo.NextFunc {
	return self.onSendToReactor
}

func (self *initParams) OnSendToConnection() rxgo.NextFunc {
	return self.onSendToConnection
}

type IConnectionReactor interface {
	io.Closer
	Init(params IInitParams) (rxgo.NextFunc, rxgo.ErrFunc, rxgo.CompletedFunc, error)
	Open() error
}

const ConnectionId = "ConnectionId"

//const UserContext = "UserContext"
