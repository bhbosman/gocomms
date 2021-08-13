package pingPong

import (
	"github.com/bhbosman/gocomms/internal"
	"reflect"
)

type StackData struct {
	outboundChannel *internal.ChannelManager
}

func NewStackData(connectionId string) *StackData {
	outboundChannel := internal.NewChannelManager("outbound PingPong", connectionId)
	return &StackData{
		outboundChannel: outboundChannel,
	}
}

func (self *StackData) Close() error {
	return nil
}

func WrongStackDataError(connectionType internal.ConnectionType, stackData interface{}) error {
	return internal.NewWrongStackDataType(
		StackName,
		connectionType,
		reflect.TypeOf((*StackData)(nil)),
		reflect.TypeOf(stackData))
}
