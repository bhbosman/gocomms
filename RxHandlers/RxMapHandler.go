package RxHandlers

import (
	"context"
	"fmt"
	model2 "github.com/bhbosman/gocommon/model"
	"github.com/bhbosman/goerrors"
	"github.com/bhbosman/goprotoextra"
	"go.uber.org/multierr"
	"go.uber.org/zap"
)

type RxMapHandler struct {
	BaseRxHandler
	next IRxMapStackHandler
}

func (self *RxMapHandler) Handler(ctx context.Context, i interface{}) (interface{}, error) {
	if rws, ok := i.(goprotoextra.ReadWriterSize); ok {
		self.RwsMessageCount++
		self.RwsByteCountIn += int64(rws.Size())
		if self.next != nil {
			newRws, err := self.next.MapReadWriterSize(ctx, rws)
			if err != nil {
				return nil, err
			}
			if newRws != nil {
				self.RwsByteCountOut += int64(newRws.Size())
				return newRws, err
			}
			return nil, fmt.Errorf("unknown error in %v, where newRws is nil", self.Name)
		}
		return i, nil
	} else {
		self.OtherMessageCount++
		if self.next != nil {
			switch v := i.(type) {
			case *model2.PublishRxHandlerCounters:
				counters := model2.NewRxHandlerCounter(
					self.Name,
					self.OtherMessageCount,
					self.RwsMessageCount,
					self.RwsByteCountIn,
					self.RwsByteCountOut)
				v.Add(counters)
				_ = self.next.ReadMessage(i)
				break
			default:
				_ = self.next.ReadMessage(i)
				break
			}
		}
		return i, nil
	}
}

func NewRxMapHandler(
	name string,
	ConnectionCancelFunc model2.ConnectionCancelFunc,
	logger *zap.Logger,
	next IRxMapStackHandler) (*RxMapHandler, error) {

	var errList error = nil
	if ConnectionCancelFunc == nil {
		errList = multierr.Append(errList, goerrors.NewInvalidNilParamError("ConnectionCancelFunc"))
	}
	if logger == nil {
		errList = multierr.Append(errList, goerrors.NewInvalidNilParamError("logger"))
	}
	if next == nil {
		errList = multierr.Append(errList, goerrors.NewInvalidNilParamError("next"))
	}
	if errList != nil {
		return nil, errList
	}
	return &RxMapHandler{
		BaseRxHandler: BaseRxHandler{
			Logger:               logger,
			RwsMessageCount:      0,
			OtherMessageCount:    0,
			RwsByteCountIn:       0,
			RwsByteCountOut:      0,
			Name:                 name,
			ConnectionCancelFunc: ConnectionCancelFunc,
		},
		next: next,
	}, nil
}
