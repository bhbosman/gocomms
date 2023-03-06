package common

import (
	"fmt"
	"github.com/bhbosman/gocommon"
	"github.com/reactivex/rxgo/v2"
	"golang.org/x/net/context"
	"strings"
)

type inboundPipeDefinition struct {
	stacks []IBoundData
}

func (self *inboundPipeDefinition) BuildIncomingObs(
	inBoundChannel chan rxgo.Item,
	stackDataMap map[string]*StackDataContainer,
	cancelCtx context.Context,
) (gocommon.IObservable, error) {
	obsIn, err := self.buildInBoundPipesObservables(stackDataMap, inBoundChannel, rxgo.WithContext(cancelCtx))
	if err != nil {
		return nil, err
	}
	return obsIn, nil
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
	for i := 0; i < len(self.stacks); i++ {
		stack := self.stacks[i].Bound()
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

func (self *inboundPipeDefinition) BuildInBoundPipeStates() ([]*PipeState, error) {
	var pipeStarts []*PipeState

	for _, currentStack := range self.stacks {
		if currentStack == nil {
			continue
		}
		stack := currentStack.Bound()
		if stack == nil {
			continue
		}
		stackBoundDefinition, err := stack()
		if err != nil {
			return nil, err
		}
		if stackBoundDefinition == nil {
			continue
		}
		pipeState := stackBoundDefinition.GetPipeState()
		if pipeState == nil {
			continue
		}
		b := true
		b = b && (0 != strings.Compare("", pipeState.ID))
		b = b && (pipeState.Create != nil)
		b = b && (pipeState.Destroy != nil)
		b = b && (pipeState.Start != nil)
		b = b && (pipeState.End != nil)
		if !b {
			return nil, fmt.Errorf("pipestate must be complete in full")
		}

		pipeStarts = append(pipeStarts, pipeState)
	}
	return pipeStarts, nil
}

type IInboundPipeDefinition interface {
	BuildInBoundPipeStates() ([]*PipeState, error)
	BuildIncomingObs(
		inBoundChannel chan rxgo.Item,
		stackDataMap map[string]*StackDataContainer,
		cancelCtx context.Context,
	) (gocommon.IObservable, error)
}

func NewInboundPipeDefinition(stacks []IBoundData) IInboundPipeDefinition {
	return &inboundPipeDefinition{stacks: stacks}
}
