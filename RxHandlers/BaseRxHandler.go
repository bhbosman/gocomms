package RxHandlers

import (
	"github.com/bhbosman/gocommon/model"
	"github.com/bhbosman/goerrors"
	"go.uber.org/multierr"
	"go.uber.org/zap"
)

type BaseRxHandler struct {
	Logger               *zap.Logger
	RwsMessageCountIn    int64
	OtherMessageCountIn  int64
	RwsMessageCountOut   int64
	OtherMessageCountOut int64
	RwsByteCountIn       int64
	RwsByteCountOut      int64
	Name                 string
	errorState           error
	ConnectionCancelFunc model.ConnectionCancelFunc
}

func NewBaseRxHandler(logger *zap.Logger, connectionCancelFunc model.ConnectionCancelFunc) (BaseRxHandler, error) {
	var errList error = nil
	if connectionCancelFunc == nil {
		errList = multierr.Append(
			errList,
			goerrors.NewInvalidNilParamError("connectionCancelFunc"),
		)
	}
	if logger == nil {
		errList = multierr.Append(errList, goerrors.NewInvalidNilParamError("logger"))
	}
	if errList != nil {
		return BaseRxHandler{}, errList
	}
	return BaseRxHandler{
		Logger:               logger,
		ConnectionCancelFunc: connectionCancelFunc,
	}, nil
}

func (self *BaseRxHandler) Close() error {
	return nil
}
