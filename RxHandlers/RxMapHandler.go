package RxHandlers

import (
	"context"
	"github.com/bhbosman/gocommon/messages"
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

	switch v := i.(type) {
	case *messages.EmptyQueue:
		return i, nil
	case *model2.PublishRxHandlerCounters:
		self.OtherMessageCountIn++
		counters := model2.NewRxHandlerCounter(
			self.Name,
			self.OtherMessageCountIn,
			self.RwsMessageCountIn,
			self.OtherMessageCountOut,
			self.RwsMessageCountOut,
			self.RwsByteCountIn,
			self.RwsByteCountOut)
		v.Add(counters)
		err := self.next.ReadMessage(i)
		if err != nil {
			return nil, err
		}
		self.OtherMessageCountOut++
		return i, nil
	default:
		switch v := i.(type) {
		case goprotoextra.ReadWriterSize:
			self.RwsMessageCountOut++
			self.RwsByteCountOut += int64(v.Size())
		default:
			self.OtherMessageCountIn++
		}
		i, err := self.next.MapReadWriterSize(ctx, i)
		if err != nil {
			return nil, err
		}
		switch v := i.(type) {
		case goprotoextra.ReadWriterSize:
			self.RwsByteCountOut += int64(v.Size())
			self.RwsMessageCountOut++
			return v, nil
		default:
			self.OtherMessageCountOut++
			err := self.next.ReadMessage(i)
			if err != nil {
				return nil, err
			}
			return i, nil
		}

	}

}

func (self *RxMapHandler) FlatMapHandler(ctx context.Context) func(item rxgo.Item) rxgo.Observable {
	return func(item rxgo.Item) rxgo.Observable {
		switch {
		case item.V != nil:
			result, err := self.next.FlatMapHandler(ctx, item.V)
			if err != nil {
				return rxgo.Just(err)()
			}
			if result.UseDefaultPath {
				unk, err := self.Handler(ctx, item.V)
				if err != nil {
					return rxgo.Just(err)()
				}
				return rxgo.Just(unk)()
			}
			if _, ok := item.V.(goprotoextra.ReadWriterSize); ok {
				self.RwsMessageCountIn++
			} else {
				self.OtherMessageCountIn++
			}
			self.RwsMessageCountOut += int64(result.RwsCount)
			self.RwsMessageCountOut += int64(result.OtherCount)
			self.RwsByteCountIn += int64(result.BytesIn)
			self.RwsByteCountOut += int64(result.BytesOut)

			return rxgo.Just(result.Items...)()
		case item.E != nil:
			return rxgo.Just(item.E)()
		default:
			return rxgo.Just()()
		}
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
			RwsMessageCountIn:    0,
			OtherMessageCountIn:  0,
			RwsMessageCountOut:   0,
			OtherMessageCountOut: 0,
			RwsByteCountIn:       0,
			RwsByteCountOut:      0,
			Name:                 name,
			ConnectionCancelFunc: ConnectionCancelFunc,
		},
		next: next,
	}, nil
}
