package common

import (
	"context"
	"fmt"
	"github.com/bhbosman/gocommon"
	"github.com/reactivex/rxgo/v2"
	"strings"
)

type pipeDefinition struct {
	stacks     []IBoundData
	isOutBound bool
}

func (self *pipeDefinition) BuildObs(
	outBoundChannel chan rxgo.Item,
	stackDataMap map[string]*StackDataContainer,
	cancelCtx context.Context,
) (gocommon.IObservable, error) {
	obsOut, err := self.buildBoundObservables(stackDataMap, outBoundChannel, rxgo.WithContext(cancelCtx))
	if err != nil {
		return nil, err
	}
	return obsOut, nil
}

func (self *pipeDefinition) buildBoundObservables(
	stackDataMap map[string]*StackDataContainer,
	outbound chan rxgo.Item,
	opts ...rxgo.Option,
) (gocommon.IObservable, error) {
	var obs gocommon.IObservable = rxgo.FromChannel(outbound, opts...)

	handleStack := func(id string, currentStack IStackBoundFactory) error {
		cb := currentStack.GetPipeDefinition()
		if cb != nil {
			var err error
			var stackData interface{} = nil
			var pipeData interface{} = nil
			if containerData, ok := stackDataMap[id]; ok {
				stackData = containerData.StackData
				pipeData =
					func(isOutBound bool) interface{} {
						if isOutBound {
							return containerData.OutPipeData
						}
						return containerData.InPipeData
					}(self.isOutBound)
			}
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
			sbd, err := stack()
			if err != nil {
				return nil, err
			}
			err = handleStack(self.stacks[i].Name(), sbd)
			if err != nil {
				return nil, err
			}
		}
	}
	return obs, nil
}

func (self *pipeDefinition) BuildPipeStates() ([]*PipeState, error) {
	var pipeStarts []*PipeState

	for _, currentStack := range self.stacks {
		if currentStack == nil {
			continue
		}
		stack := currentStack.Bound()
		if stack == nil {
			continue
		}
		sbd, err := stack()
		if err != nil {
			return nil, err
		}
		if sbd == nil {
			continue
		}
		pipeState := sbd.GetPipeState()
		if pipeState == nil {
			continue
		}
		b := true
		b = b && (0 != strings.Compare("", pipeState.ID))
		b = b && (pipeState.OnCreate != nil)
		b = b && (pipeState.OnDestroy != nil)
		b = b && (pipeState.OnStart != nil)
		b = b && (pipeState.OnEnd != nil)
		if !b {
			return nil, fmt.Errorf("stackstate must be complete in full")
		}
		pipeStarts = append(pipeStarts, pipeState)
	}

	return pipeStarts, nil
}

type IPipeDefinition interface {
	BuildPipeStates() ([]*PipeState, error)
	BuildObs(
		outBoundChannel chan rxgo.Item,
		stackDataMap map[string]*StackDataContainer,
		cancelCtx context.Context,
	) (gocommon.IObservable, error)
}

func NewPipeDefinition(stacks []IBoundData, isOutBound bool) IPipeDefinition {
	return &pipeDefinition{
		stacks:     stacks,
		isOutBound: isOutBound,
	}
}
