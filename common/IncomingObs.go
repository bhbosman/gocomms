package common

import (
	"github.com/bhbosman/gocommon"
)

type IncomingObs struct {
	InboundObservable gocommon.IObservable
}

func NewIncomingObs(inboundObservable gocommon.IObservable) *IncomingObs {
	return &IncomingObs{
		InboundObservable: inboundObservable,
	}
}
