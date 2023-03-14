package RxHandlers

import (
	"context"
	"github.com/bhbosman/gocommon/messages"
	"github.com/bhbosman/gocommon/model"
	"github.com/bhbosman/goerrors"
	"github.com/bhbosman/goprotoextra"
	rxgo "github.com/reactivex/rxgo/v2"
	"go.uber.org/multierr"
	"go.uber.org/zap"
)

type RxMapHandler struct {
	BaseRxHandler
	next IRxMapStackHandler
}

func (self *RxMapHandler) Handler(ctx context.Context, unk interface{}) (interface{}, error) {
	// TODO: need to clean up this bit, on how read message is called and how it affect RWS, with protobuf, and clear counters. For now it is hard coded
	switch v01 := unk.(type) {
	case goprotoextra.ReadWriterSize:
		self.RwsMessageCountIn++
		self.RwsByteCountIn += int64(v01.Size())
		ans, err := self.next.MapReadWriterSize(ctx, v01)
		if err != nil {
			return nil, err
		}
		switch v02 := ans.(type) {
		case goprotoextra.ReadWriterSize:
			self.RwsByteCountOut += int64(v02.Size())
			self.RwsMessageCountOut++
			return v02, nil
		default:
			self.OtherMessageCountOut++
			return v02, nil
		}
	case *model.ClearCounters:
		self.clearCounters()
		self.next.ClearCounters()
		return v01, nil
	case *messages.EmptyQueue:
		self.next.EmptyQueue()
		return v01, nil
	case *model.PublishRxHandlerCounters:
		self.OtherMessageCountIn++
		counters := model.NewRxHandlerCounter(
			self.Name,
			self.OtherMessageCountIn,
			self.RwsMessageCountIn,
			self.OtherMessageCountOut,
			self.RwsMessageCountOut,
			self.RwsByteCountIn,
			self.RwsByteCountOut)
		v01.Add(counters)
		self.next.PublishCounters(v01)
		self.OtherMessageCountOut++
		return v01, nil
	default:
		self.OtherMessageCountIn++
		ans, err := self.next.MapReadWriterSize(ctx, v01)
		if err != nil {
			return nil, err
		}
		switch v02 := ans.(type) {
		case goprotoextra.ReadWriterSize:
			self.RwsByteCountOut += int64(v02.Size())
			self.RwsMessageCountOut++
			return v02, nil
		default:
			self.OtherMessageCountOut++
			return v02, nil
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
	ConnectionCancelFunc model.ConnectionCancelFunc,
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
