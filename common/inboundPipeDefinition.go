package common

import (
	"github.com/bhbosman/gocommon"
	"github.com/reactivex/rxgo/v2"
	"golang.org/x/net/context"
)

type inboundPipeDefinition struct {
	stacks []IStackDefinition
}

func (self *inboundPipeDefinition) BuildIncomingObs(
	inBoundChannel chan rxgo.Item,
	stackDataMap map[string]*StackDataContainer,
	cancelCtx context.Context,
) (*IncomingObs, error) {
	obsIn, err := self.buildInBoundPipesObservables(stackDataMap, inBoundChannel, rxgo.WithContext(cancelCtx))
	if err != nil {
		return nil, err
	}
	return NewIncomingObs(obsIn), nil
}

func (self *inboundPipeDefinition) buildInBoundPipesObservables(
	stackDataMap map[string]*StackDataContainer,
	inbound chan rxgo.Item,
	opts ...rxgo.Option,
) (gocommon.IObservable, error) {
	var obs gocommon.IObservable = rxgo.FromChannel(inbound, opts...)

	handleStack := func(id string, currentStack IStackBoundDefinition) error {
		cb := currentStack.GetPipeDefinition()
		if cb != nil {
			var stackData interface{} = nil
			var pipeData interface{} = nil
			if container, ok := stackDataMap[id]; ok {
				stackData = container.StackData
				pipeData = container.InPipeData
			}

			var err error
			obs, err = cb(stackData, pipeData, obs)
			if err != nil {
				return err
			}
		}
		return nil
	}
	for i := len(self.stacks) - 1; i >= 0; i-- {
		stack := self.stacks[i].Inbound()
		if stack != nil {
			stackBoundDefinition, err := stack()
			if err != nil {
				return nil, err
			}
			err = handleStack(self.stacks[i].Name(), stackBoundDefinition)
			if err != nil {
				return nil, err
			}
		}
	}
	return obs, nil
}
