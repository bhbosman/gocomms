package common

import (
	"fmt"
	"github.com/bhbosman/gocommon"
	"github.com/reactivex/rxgo/v2"
	"golang.org/x/net/context"
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
	var err error
	var obsOut gocommon.IObservable
	obsOut, err = self.buildOutBoundObservables(stackDataMap, outBoundChannel, rxgo.WithContext(cancelCtx))
	if err != nil {
		return nil, err
	}

	return obsOut, nil
}

func (self *pipeDefinition) buildOutBoundObservables(
	stackDataMap map[string]*StackDataContainer,
	outbound chan rxgo.Item,
	opts ...rxgo.Option,
) (gocommon.IObservable, error) {
	var obs gocommon.IObservable = rxgo.FromChannel(outbound, opts...)

	handleStack := func(id string, currentStack IStackBoundDefinition) error {
		cb := currentStack.GetPipeDefinition()
		if cb != nil {
			var err error
			var stackData interface{} = nil
			var pipeData interface{} = nil
			if containerData, ok := stackDataMap[id]; ok {
				stackData = containerData.StackData
				pipeData = func(isOutBound bool) interface{} {
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

func (self *pipeDefinition) BuildOutBoundPipeStates() ([]*PipeState, error) {
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
			return nil, fmt.Errorf("stackstate must be complete in full")
		}
		pipeStarts = append(pipeStarts, pipeState)
	}

	return pipeStarts, nil
}

type IPipeDefinition interface {
	BuildOutBoundPipeStates() ([]*PipeState, error)
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
