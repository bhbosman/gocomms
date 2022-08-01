package RxHandlers

import (
	"context"
	model2 "github.com/bhbosman/gocommon/model"
	"github.com/bhbosman/goerrors"
	"github.com/bhbosman/goprotoextra"
	"github.com/reactivex/rxgo/v2"
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

		newRws, err := self.next.MapReadWriterSize(ctx, rws)
		if err != nil {
			return nil, err
		}
		message, b, err := self.next.ReadMessage(rws)
		if err != nil {
			return nil, err
		}
		if b {
			self.OtherMessageCount++
			return message, nil
		}

		self.RwsByteCountOut += int64(newRws.Size())
		return newRws, err
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
				_, _, err := self.next.ReadMessage(i)
				if err != nil {
					return nil, err
				}
				break
			default:
				message, b, err := self.next.ReadMessage(i)
				if err != nil {
					return nil, err
				}
				if rws, ok := message.(goprotoextra.ReadWriterSize); ok {
					self.RwsByteCountOut += int64(rws.Size())
				}
				if b {
					return message, nil
				}
				break
			}
		}
		return i, nil
	}
}

func (self *RxMapHandler) FlatMapHandler(item rxgo.Item) rxgo.Observable {
	return self.next.FlatMapHandler(item)
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
