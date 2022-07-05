package common

import (
	"github.com/bhbosman/goConnectionManager"
	"github.com/bhbosman/gocommon/model"
	"github.com/bhbosman/gocomms/RxHandlers"
	"github.com/bhbosman/goprotoextra"
)

type InvokeInboundTransportLayerHandler struct {
	isDisposed bool
	errorState error
	conn       goConnectionManager.IPublishConnectionInformation
	sendData   func(data interface{})
	sendOther  func(interface{}) error

	sendError func(err error)
	complete  func()
}

func (self *InvokeInboundTransportLayerHandler) GetAdditionalBytesIncoming() int {
	return 0
}

func (self *InvokeInboundTransportLayerHandler) SendError(err error) {
	if self.sendError != nil {
		self.sendError(err)
	}
}

func (self *InvokeInboundTransportLayerHandler) GetAdditionalBytesSend() int {
	return 0
}

func (self *InvokeInboundTransportLayerHandler) ReadMessage(i interface{}) error {
	if publishRxHandlerCounters, ok := i.(*model.PublishRxHandlerCounters); ok {
		return self.conn.ConnectionInformationReceived(publishRxHandlerCounters)
	}
	return nil
}

func NewInvokeInboundTransportLayerHandler(
	conn goConnectionManager.IPublishConnectionInformation,
	sendData func(data interface{}),
	sendOther func(interface{}) error,
	sendError func(err error),
	complete func(),
) (*InvokeInboundTransportLayerHandler, error) {
	return &InvokeInboundTransportLayerHandler{
		errorState: nil,
		conn:       conn,
		sendData:   sendData,
		sendOther:  sendOther,
		sendError:  sendError,
		complete:   complete,
	}, nil
}

func (self *InvokeInboundTransportLayerHandler) Close() error {
	if self.isDisposed {
		return self.errorState
	}
	self.isDisposed = true
	return nil
}

func (self *InvokeInboundTransportLayerHandler) SendData(data interface{}) {
	_ = self.ReadMessage(data)
	if self.sendData != nil {
		self.sendData(data)
	}
}

func (self *InvokeInboundTransportLayerHandler) OnError(err error) {
	self.errorState = err
}

func (self *InvokeInboundTransportLayerHandler) NextReadWriterSize(
	rws goprotoextra.ReadWriterSize,
	_ func(rws goprotoextra.ReadWriterSize) error,
	_ func(interface{}) error,
	sizeUpdate func(size int) error) error {

	if self.errorState != nil {
		return self.errorState
	}
	_ = sizeUpdate(rws.Size())
	if self.sendData != nil {
		self.sendData(
			rws,
			//RxHandlers.NewNextExternal(true, rws),
		)
	}
	return self.errorState
}

func (self *InvokeInboundTransportLayerHandler) OnComplete() {
	if self.errorState == nil {
		self.errorState = RxHandlers.RxHandlerComplete
	}
}

func (self *InvokeInboundTransportLayerHandler) Complete() {
	if self.complete != nil {
		self.complete()
	}
}
