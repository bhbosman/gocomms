package RxHandlers

import (
	"github.com/bhbosman/goCommsDefinitions"
	"github.com/bhbosman/gocommon/messages"
	model2 "github.com/bhbosman/gocommon/model"
	"github.com/bhbosman/goerrors"
	"github.com/bhbosman/goprotoextra"
	"github.com/reactivex/rxgo/v2"
	"go.uber.org/multierr"
	"go.uber.org/zap"
)

type RxNextHandler struct {
	BaseRxHandler
	stackHandler  IRxNextStackHandler
	onSendData    rxgo.NextFunc
	onSendError   rxgo.ErrFunc
	onComplete    rxgo.CompletedFunc
	onTrySendData goCommsDefinitions.TryNextFunc
	isActive      func() bool
}

func (self *RxNextHandler) IsActive() bool {
	return self.isActive()
}

func (self *RxNextHandler) OnTrySendData(i interface{}) bool {
	return self.onTrySendData(i)
}

func (self *RxNextHandler) Close() error {
	var err error
	err = multierr.Append(err, self.BaseRxHandler.Close())
	return err
}

func (self *RxNextHandler) toNextChannelWithByteOutCount(rws goprotoextra.ReadWriterSize) error {
	rwsSize := rws.Size()
	if rwsSize != 0 {
		self.RwsMessageCountOut++
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
	self.RwsMessageCountOut++
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
		self.RwsMessageCountIn++
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
		switch v := i.(type) {
		case *messages.EmptyQueue:
			break
		case *model2.PublishRxHandlerCounters:
			self.OtherMessageCountIn++
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
				self.OtherMessageCountIn,
				self.RwsMessageCountIn,
				self.OtherMessageCountOut,
				self.RwsMessageCountOut,
				bytesInCount,
				byteOutCount)
			v.Add(counter)
			if self.stackHandler != nil {
				self.stackHandler.ReadMessage(i)
			}
			if self.onSendData != nil {
				self.OtherMessageCountOut++
				self.onSendData(i)
			}
			break
		default:
			self.OtherMessageCountIn++
			if self.stackHandler != nil {
				self.stackHandler.ReadMessage(i)
			}
			if self.onSendData != nil {
				self.OtherMessageCountOut++
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
func newRxNextHandler(
	name string,
	ConnectionCancelFunc model2.ConnectionCancelFunc,
	next IRxNextStackHandler,
	onSendData rxgo.NextFunc,
	onTrySendData goCommsDefinitions.TryNextFunc,
	onSendError rxgo.ErrFunc,
	onComplete rxgo.CompletedFunc,
	isActive func() bool,
	logger *zap.Logger) (*RxNextHandler, error) {

	if onSendData == nil {
		return nil, goerrors.InvalidParam
	}
	if onTrySendData == nil {
		return nil, goerrors.InvalidParam
	}
	if onSendError == nil {
		return nil, goerrors.InvalidParam
	}
	if onComplete == nil {
		return nil, goerrors.InvalidParam
	}

	return &RxNextHandler{
		BaseRxHandler: BaseRxHandler{
			Logger:               logger,
			RwsMessageCountIn:    0,
			OtherMessageCountIn:  0,
			RwsMessageCountOut:   0,
			OtherMessageCountOut: 0,
			RwsByteCountIn:       0,
			RwsByteCountOut:      0,
			Name:                 name,
			errorState:           nil,
			ConnectionCancelFunc: ConnectionCancelFunc,
		},
		stackHandler:  next,
		onSendData:    onSendData,
		onTrySendData: onTrySendData,
		onSendError:   onSendError,
		onComplete:    onComplete,
		isActive: func(isActive func() bool) func() bool {
			if isActive != nil {
				return isActive
			}
			return func() bool {
				return true
			}
		}(isActive),
	}, nil
}

func NewRxNextHandler2(
	name string,
	ConnectionCancelFunc model2.ConnectionCancelFunc,
	next IRxNextStackHandler,
	defaultRxNextHandler goCommsDefinitions.IRxNextHandler,
	logger *zap.Logger,
) (*RxNextHandler, error) {
	if defaultRxNextHandler == nil {
		return nil, goerrors.InvalidParam
	}
	return newRxNextHandler(
		name,
		ConnectionCancelFunc,
		next,
		defaultRxNextHandler.OnSendData,
		defaultRxNextHandler.OnTrySendData,
		defaultRxNextHandler.OnError,
		defaultRxNextHandler.OnComplete,
		defaultRxNextHandler.IsActive,
		logger)
}
