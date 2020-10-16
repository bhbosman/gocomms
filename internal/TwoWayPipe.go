package internal

import (
	"context"
	"github.com/bhbosman/gologging"
	"github.com/bhbosman/goprotoextra"
	"github.com/reactivex/rxgo/v2"
	"go.uber.org/multierr"
	"io"
)

type TwoWayPipe struct {
	logger                 *gologging.SubSystemLogger
	InboundObservable      rxgo.Observable
	OutboundObservable     rxgo.Observable
	cancelCtx              context.Context
	PipeState              []PipeState
	StackState             []StackState
	inboundChannelManager  *ChannelManager
	outboundChannelManager *ChannelManager
}

func (self *TwoWayPipe) SendOutgoingData(rws goprotoextra.ReadWriterSize) error {
	err := self.cancelCtx.Err()
	if err != nil {
		return err
	}
	self.outboundChannelManager.Send(self.cancelCtx, rws)
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
	self.inboundChannelManager.Send(self.cancelCtx, item)

	return nil
}

func (self *TwoWayPipe) SendError(item error) error {
	err := self.cancelCtx.Err()
	if err != nil {
		return err
	}
	self.outboundChannelManager.SendError(self.cancelCtx, item)
	return nil
}

func (self *TwoWayPipe) Close() error {
	var err error
	err = multierr.Append(err, self.outboundChannelManager.Close())
	err = multierr.Append(err, self.inboundChannelManager.Close())
	return err
}

func NewTwoWayPipe(
	connectionId string,
	logger *gologging.SubSystemLogger,
	InBound chan rxgo.Item,
	OutBound chan rxgo.Item,
	InboundObservable rxgo.Observable,
	OutboundObservable rxgo.Observable,
	cancelCtx context.Context,
	pipeStarts []PipeState,
	stackState []StackState) *TwoWayPipe {
	return &TwoWayPipe{
		inboundChannelManager:  NewChannelManager(InBound, "inboundChannelManager", connectionId),
		outboundChannelManager: NewChannelManager(OutBound, "outboundChannelManager", connectionId),
		logger:                 logger,
		InboundObservable:      InboundObservable,
		OutboundObservable:     OutboundObservable,
		cancelCtx:              cancelCtx,
		PipeState:              pipeStarts,
		StackState:             stackState,
	}
}
