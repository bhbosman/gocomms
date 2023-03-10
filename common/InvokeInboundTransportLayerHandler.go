package common

import (
	"github.com/bhbosman/gocomms/RxHandlers"
	"github.com/bhbosman/goprotoextra"
)

type InvokeInboundTransportLayerHandler struct {
	isDisposed  bool
	errorState  error
	sendData    func(data interface{})
	trySendData func(data interface{}) bool
	sendError   func(err error)
	complete    func()
}

func (self *InvokeInboundTransportLayerHandler) IsActive() bool {
	return true
}

func (self *InvokeInboundTransportLayerHandler) OnSendData(i interface{}) {
	self.sendData(i)
}

func (self *InvokeInboundTransportLayerHandler) OnTrySendData(i interface{}) bool {
	return self.trySendData(i)
}

func (self *InvokeInboundTransportLayerHandler) GetAdditionalBytesIncoming() int {
	return 0
}

func (self *InvokeInboundTransportLayerHandler) GetAdditionalBytesSend() int {
	return 0
}

func (self *InvokeInboundTransportLayerHandler) ReadMessage(i interface{}) error {
	return nil
}

func (self *InvokeInboundTransportLayerHandler) Close() error {
	if self.isDisposed {
		return self.errorState
	}
	self.isDisposed = true
	return nil
}

func (self *InvokeInboundTransportLayerHandler) OnError(err error) {
	self.errorState = err
}

func (self *InvokeInboundTransportLayerHandler) NextReadWriterSize(
	rws goprotoextra.ReadWriterSize,
	_ func(rws goprotoextra.ReadWriterSize) error,
	_ func(interface{}) error,
	_ func(size int) error,
) error {
	if self.errorState != nil {
		return self.errorState
	}
	if self.sendData != nil {
		self.sendData(rws)
	}
	return self.errorState
}

func (self *InvokeInboundTransportLayerHandler) OnComplete() {
	if self.errorState == nil {
		self.errorState = RxHandlers.RxHandlerComplete
	}
}

func NewInvokeInboundTransportLayerHandler(
	sendData func(data interface{}),
	trySendData func(data interface{}) bool,
	sendError func(err error),
	complete func(),
) (*InvokeInboundTransportLayerHandler, error) {
	return &InvokeInboundTransportLayerHandler{
		errorState:  nil,
		sendData:    sendData,
		trySendData: trySendData,
		sendError:   sendError,
		complete:    complete,
	}, nil
}
