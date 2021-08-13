package internal

import (
	"context"
	"github.com/reactivex/rxgo/v2"
	"go.uber.org/multierr"
	"io"
)

type IncomingObs struct {
	InboundObservable     rxgo.Observable
	InboundChannelManager *ChannelManager
	CancelCtx             context.Context
}

func (self *IncomingObs) Close() error {
	var err error = nil
	err = multierr.Append(err, self.InboundChannelManager.Close())
	return err
}

func (self *IncomingObs) ReceiveIncomingData(item io.Reader) error {
	err := self.CancelCtx.Err()
	if err != nil {
		return err
	}
	self.InboundChannelManager.Send(self.CancelCtx, item)
	return nil
}
