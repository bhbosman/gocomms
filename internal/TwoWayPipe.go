package internal

import (
	"context"
	"github.com/bhbosman/gologging"
	"github.com/bhbosman/goprotoextra"
	"github.com/reactivex/rxgo/v2"
	"io"
)

type TwoWayPipe struct {
	logger             *gologging.SubSystemLogger
	InBound            chan rxgo.Item
	OutBound           chan rxgo.Item
	InboundObservable  rxgo.Observable
	OutboundObservable rxgo.Observable
	cancelCtx          context.Context
	PipeState          []PipeState
	StackState         []StackState
}

func (self *TwoWayPipe) SendOutgoingData(rws goprotoextra.ReadWriterSize) error {
	err := self.cancelCtx.Err()
	if err != nil {
		return err
	}
	item := rxgo.Of(rws)
	item.SendContext(self.cancelCtx, self.OutBound)
	return nil
}

func (self *TwoWayPipe) ReceiveIncomingData(item io.Reader) error {
	err := self.cancelCtx.Err()
	if err != nil {
		return err
	}
	defer func() {
		r := recover()
		if r != nil {
			if err, ok := r.(error); ok {
				_ = self.logger.ErrorWithDescription("ReceiveIncomingData", err)
			}
		}
	}()
	self.InBound <- rxgo.Of(item)
	r := recover()
	if err, ok := r.(error); ok {
		return err
	}
	return nil
}

func (self *TwoWayPipe) SendError(item error) error {
	err := self.cancelCtx.Err()
	if err != nil {
		return err
	}
	self.OutBound <- rxgo.Error(item)
	return nil
}

func (self *TwoWayPipe) Close() error {
	close(self.OutBound)
	close(self.InBound)
	return nil
}

func NewTwoWayPipe(
	logger *gologging.SubSystemLogger,
	InBound chan rxgo.Item,
	OutBound chan rxgo.Item,
	InboundObservable rxgo.Observable,
	OutboundObservable rxgo.Observable,
	cancelCtx context.Context,
	pipeStarts []PipeState,
	stackState []StackState) *TwoWayPipe {
	return &TwoWayPipe{
		logger:             logger,
		InBound:            InBound,
		OutBound:           OutBound,
		InboundObservable:  InboundObservable,
		OutboundObservable: OutboundObservable,
		cancelCtx:          cancelCtx,
		PipeState:          pipeStarts,
		StackState:         stackState,
	}
}
