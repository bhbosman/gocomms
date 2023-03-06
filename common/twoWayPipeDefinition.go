package common

import (
	"context"
	"github.com/reactivex/rxgo/v2"
)

type IOutboundPipeDefinition interface {
	BuildOutBoundPipeStates() ([]*PipeState, error)
	BuildOutgoingObs(
		outBoundChannel chan rxgo.Item,
		stackDataMap map[string]*StackDataContainer,
		cancelCtx context.Context,
	) (*OutgoingObs, error)
}
