package RxHandlers

import (
	model2 "github.com/bhbosman/gocommon/model"
	"github.com/bhbosman/goprotoextra"
	"github.com/reactivex/rxgo/v2"
	"go.uber.org/multierr"
	"go.uber.org/zap"
)

type RxNextHandler struct {
	BaseRxHandler
	stackHandler IRxNextStackHandler
	onSendData   rxgo.NextFunc
	onSendError  rxgo.ErrFunc
	onComplete   rxgo.CompletedFunc
}

func (self *RxNextHandler) Close() error {
	var err error
	err = multierr.Append(err, self.BaseRxHandler.Close())
	return err
}

func (self *RxNextHandler) toNextChannelWithByteOutCount(rws goprotoextra.ReadWriterSize) error {
	rwsSize := rws.Size()
	if rwsSize != 0 {
		self.RwsByteCountOut += int64(rws.Size())
		if self.onSendData != nil {
			self.onSendData(rws)
		}
	}
	return nil
}

func (self *RxNextHandler) updateOutByteCount(size int) error {
	if self.errorState != nil {
		return self.errorState
	}
	self.RwsByteCountOut += int64(size)
	return nil
}

func (self *RxNextHandler) OnSendData(i interface{}) {
	// What ever code is added there, please adhere to two states:
	// 		self.stackHandler assigned and self.net not assigned
	if self.errorState != nil {
		return
	}
	var ok bool
	var rws goprotoextra.ReadWriterSize
	if rws, ok = i.(goprotoextra.ReadWriterSize); ok {
		self.RwsMessageCount++
		self.RwsByteCountIn += int64(rws.Size())
		if self.stackHandler != nil {
			err := self.stackHandler.NextReadWriterSize(
				rws,
				self.toNextChannelWithByteOutCount,
				func(i interface{}) error {
					if self.onSendData != nil {
						self.onSendData(i)
					}
					return nil
				},
				self.updateOutByteCount)
			if err != nil {
				self.setErrorState(err, true)
			}
		} else {
			err := self.toNextChannelWithByteOutCount(rws)
			if err != nil {
				return
			}
		}
	} else {
		self.OtherMessageCount++
		switch v := i.(type) {
		case *model2.PublishRxHandlerCounters:
			byteOutCount := self.RwsByteCountOut
			if self.stackHandler != nil {
				byteOutCount += int64(self.stackHandler.GetAdditionalBytesSend())
			}
			bytesInCount := self.RwsByteCountIn
			if self.stackHandler != nil {
				bytesInCount += int64(self.stackHandler.GetAdditionalBytesIncoming())
			}
			counter := model2.NewRxHandlerCounter(
				self.Name,
				self.OtherMessageCount,
				self.RwsMessageCount,
				bytesInCount,
				byteOutCount)
			v.Add(counter)
			if self.stackHandler != nil {
				self.stackHandler.ReadMessage(i)
			}
			if self.onSendData != nil {
				self.onSendData(i)
			}
			break
		default:
			if self.stackHandler != nil {
				self.stackHandler.ReadMessage(i)
			}
			if self.onSendData != nil {
				self.onSendData(i)
			}
			break
		}
	}
}

func (self *RxNextHandler) OnError(err error) {
	self.setErrorState(err, false)
	if self.stackHandler != nil {
		self.stackHandler.OnError(err)
	}
	if self.onSendError != nil {
		self.onSendError(err)
	}
}

func (self *RxNextHandler) OnComplete() {
	self.setErrorState(RxHandlerComplete, false)
	if self.onComplete != nil {
		self.onComplete()
	}
}

func (self *RxNextHandler) setErrorState(err error, propagateError bool) {
	if self.errorState == nil {
		self.errorState = err
		if propagateError {
			self.OnError(err)
		}
	}
}

func NewRxNextHandler(
	name string,
	ConnectionCancelFunc model2.ConnectionCancelFunc,
	next IRxNextStackHandler,
	onSendData func(data interface{}),
	onSendError func(err error),
	onComplete func(),
	logger *zap.Logger) (*RxNextHandler, error) {
	return &RxNextHandler{
		BaseRxHandler: BaseRxHandler{
			Logger:               logger,
			RwsMessageCount:      0,
			OtherMessageCount:    0,
			RwsByteCountIn:       0,
			RwsByteCountOut:      0,
			Name:                 name,
			errorState:           nil,
			ConnectionCancelFunc: ConnectionCancelFunc,
		},
		stackHandler: next,
		onSendData:   onSendData,
		onSendError:  onSendError,
		onComplete:   onComplete,
	}, nil
}
