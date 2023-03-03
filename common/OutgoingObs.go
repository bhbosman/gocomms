package common

import (
	"context"
	"github.com/bhbosman/gocommon"
)

type OutgoingObs struct {
	OutboundObservable gocommon.IObservable
	CancelCtx          context.Context
}

func NewOutgoingObs(
	outboundObservable gocommon.IObservable,
	cancelCtx context.Context,
) *OutgoingObs {
	return &OutgoingObs{
		OutboundObservable: outboundObservable,
		CancelCtx:          cancelCtx}
}
