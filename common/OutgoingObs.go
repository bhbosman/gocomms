package common

import (
	"context"
	"github.com/reactivex/rxgo/v2"
)

type OutgoingObs struct {
	OutboundObservable rxgo.Observable
	CancelCtx          context.Context
}

func NewOutgoingObs(outboundObservable rxgo.Observable, cancelCtx context.Context) *OutgoingObs {
	return &OutgoingObs{
		OutboundObservable: outboundObservable,
		CancelCtx:          cancelCtx}
}
