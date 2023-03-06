package common

import (
	"github.com/bhbosman/gocommon"
)

type OutgoingObs struct {
	OutboundObservable gocommon.IObservable
	//CancelCtx          context.Context
}

func NewOutgoingObs(
	outboundObservable gocommon.IObservable,
) *OutgoingObs {
	return &OutgoingObs{
		OutboundObservable: outboundObservable,
		//CancelCtx:          cancelCtx,
	}
}
