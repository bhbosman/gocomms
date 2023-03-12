package intf

import "github.com/reactivex/rxgo/v2"

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

func (self *initParams) OnSendToReactor() rxgo.NextFunc {
	return self.onSendToReactor
}

func (self *initParams) OnSendToConnection() rxgo.NextFunc {
	return self.onSendToConnection
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
