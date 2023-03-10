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

type initParams struct {
	onSendToReactor         rxgo.NextFunc
	onSendToConnection      rxgo.NextFunc
	nextFuncOutBoundChannel rxgo.NextFunc
	nextFuncInBoundChannel  rxgo.NextFunc
}

func (self *initParams) NextFuncOutBoundChannel() rxgo.NextFunc {
	return self.nextFuncOutBoundChannel
}

func (self *initParams) NextFuncInBoundChannel() rxgo.NextFunc {
	return self.nextFuncInBoundChannel
}

func NewInitParams(
	onSendToReactor rxgo.NextFunc,
	onSendToConnection rxgo.NextFunc,
	nextFuncOutBoundChannel rxgo.NextFunc,
	nextFuncInBoundChannel rxgo.NextFunc,
) IInitParams {
	return &initParams{
		onSendToReactor:         onSendToReactor,
		onSendToConnection:      onSendToConnection,
		nextFuncOutBoundChannel: nextFuncOutBoundChannel,
		nextFuncInBoundChannel:  nextFuncInBoundChannel,
	}
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
