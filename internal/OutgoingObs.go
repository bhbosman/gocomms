package internal

import (
	"context"
	"github.com/bhbosman/goprotoextra"
	"github.com/reactivex/rxgo/v2"
	"go.uber.org/multierr"
)

type OutgoingObs struct {
	OutboundChannelManager *ChannelManager
	OutboundObservable     rxgo.Observable
	CancelCtx              context.Context
}

func (self *OutgoingObs) Close() error {
	var err error = nil
	err = multierr.Append(err, self.OutboundChannelManager.Close())
	return err
}

func (self *OutgoingObs) SendOutgoingData(rws goprotoextra.ReadWriterSize) error {
	err := self.CancelCtx.Err()
	if err != nil {
		return err
	}
	self.OutboundChannelManager.Send(self.CancelCtx, rws)
	return nil
}
