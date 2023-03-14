package common

import (
	"github.com/bhbosman/goConnectionManager"
	"github.com/bhbosman/gocommon/model"
	"github.com/bhbosman/gocomms/RxHandlers"
	"github.com/bhbosman/goprotoextra"
	"io"
)

type InvokeWriterHandler struct {
	errorState error
	Writer     io.Writer
	oci        goConnectionManager.IPublishConnectionInformation
}

func (self *InvokeWriterHandler) PublishCounters(counters *model.PublishRxHandlerCounters) {
	_ = self.oci.ConnectionInformationReceived(counters)
}

func (self *InvokeWriterHandler) EmptyQueue() {
}

func (self *InvokeWriterHandler) ClearCounters() {
}

func (self *InvokeWriterHandler) IsActive() bool {
	return true
}

func (self *InvokeWriterHandler) OnSendData(i interface{}) {
	if self.errorState != nil {
		return
	}

}

func (self *InvokeWriterHandler) OnTrySendData(i interface{}) bool {
	if self.errorState != nil {
		return false
	}

	return true
}

func (self *InvokeWriterHandler) GetAdditionalBytesIncoming() int {
	return 0
}

func (self *InvokeWriterHandler) SendError(_ error) {
}

func (self *InvokeWriterHandler) GetAdditionalBytesSend() int {
	return 0
}

func (self *InvokeWriterHandler) Close() error {
	return nil
}

func (self *InvokeWriterHandler) OnError(err error) {
	if self.errorState == nil {
		self.errorState = err
	}
}

func (self *InvokeWriterHandler) NextReadWriterSize(
	rws goprotoextra.ReadWriterSize,
	_ func(rws goprotoextra.ReadWriterSize) error,
	_ func(interface{}) error,
	sizeUpdate func(size int) error) error {

	if self.errorState != nil {
		return self.errorState
	}
	_ = sizeUpdate(rws.Size())
	_, err := io.Copy(self.Writer, rws)
	if err != nil {
		self.errorState = err
	}
	return self.errorState
}

func (self *InvokeWriterHandler) OnComplete() {
	if self.errorState == nil {
		self.errorState = RxHandlers.RxHandlerComplete
	}
}

func (self *InvokeWriterHandler) Complete() {

}

func NewInvokeOutBoundTransportLayerHandler(
	writer io.Writer,
	pci goConnectionManager.IPublishConnectionInformation) *InvokeWriterHandler {
	return &InvokeWriterHandler{
		errorState: nil,
		Writer:     writer,
		oci:        pci,
	}
}
