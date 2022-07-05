package common

import (
	"github.com/reactivex/rxgo/v2"
)

type IncomingObs struct {
	InboundObservable rxgo.Observable
}

func NewIncomingObs(inboundObservable rxgo.Observable) *IncomingObs {
	return &IncomingObs{
		InboundObservable: inboundObservable,
	}
}
