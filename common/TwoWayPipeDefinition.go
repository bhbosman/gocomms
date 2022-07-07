package common

import (
	"context"
	"fmt"
	"github.com/bhbosman/goerrors"
	"github.com/reactivex/rxgo/v2"
	"strings"
)

type IStackDefinitionIndex interface {
	GetStackDefinition() IStackDefinition
}

type StackDefinitionIndex struct {
	StackDefinition IStackDefinition
}

func (self *StackDefinitionIndex) GetStackDefinition() IStackDefinition {
	return self.StackDefinition
}

type TwoWayPipeDefinition struct {
	Stacks []IStackDefinitionIndex
}

func (self *TwoWayPipeDefinition) AddStackDefinition(stack IStackDefinition) {
	self.Stacks = append(
		self.Stacks,
		&StackDefinitionIndex{
			StackDefinition: stack,
		})
}

func (self TwoWayPipeDefinition) BuildStackState() ([]*StackState, error) {
	var allStackState []*StackState
	for _, item := range self.Stacks {
		stackState := item.GetStackDefinition().GetStackState()
		if stackState == nil {
			continue
		}
		b := true
		b = b && (0 != strings.Compare("", stackState.Id))
		b = b && (stackState.Create != nil)
		b = b && (stackState.Destroy != nil)
		b = b && (stackState.Start != nil)
		b = b && (stackState.Stop != nil)
		if !b {
			return nil, fmt.Errorf("stackstate must be complete in full")
		}
		allStackState = append(allStackState, stackState)

	}
	return allStackState, nil
}

func (self TwoWayPipeDefinition) BuildIncomingObs(inBoundChannel chan rxgo.Item, stackDataMap map[string]*StackDataContainer, cancelCtx context.Context) (*IncomingObs, error) {
	obsIn, err := self.buildInBoundPipesObservables(stackDataMap, inBoundChannel, rxgo.WithContext(cancelCtx))
	if err != nil {
		return nil, err
	}
	return NewIncomingObs(obsIn), nil
}

func (self TwoWayPipeDefinition) BuildOutgoingObs(
	outBoundChannel chan rxgo.Item,
	stackDataMap map[string]*StackDataContainer,
	cancelCtx context.Context) (*OutgoingObs, error) {

	var err error
	var obsOut rxgo.Observable
	obsOut, err = self.buildOutBoundObservables(stackDataMap, outBoundChannel, rxgo.WithContext(cancelCtx))
	if err != nil {
		return nil, err
	}

	return NewOutgoingObs(obsOut, cancelCtx), nil
}

func (self *TwoWayPipeDefinition) BuildOutBoundPipeStates() ([]*PipeState, error) {
	var pipeStarts []*PipeState

	for _, currentStack := range self.Stacks {
		if currentStack == nil {
			continue
		}
		stack := currentStack.GetStackDefinition().GetOutbound()
		if stack == nil {
			continue
		}
		boundResult, err := stack.GetBoundResult()
		if err != nil {
			return nil, err
		}
		if boundResult == nil {
			continue
		}
		var stackBoundDefinition IStackBoundDefinition
		stackBoundDefinition, err = boundResult()
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

func (self *TwoWayPipeDefinition) buildOutBoundObservables(stackDataMap map[string]*StackDataContainer, outbound chan rxgo.Item, opts ...rxgo.Option) (rxgo.Observable, error) {
	obs := rxgo.FromChannel(outbound, opts...)

	handleStack := func(id string, currentStack IStackBoundDefinition) error {
		cb := currentStack.GetPipeDefinition()
		if cb != nil {
			var err error
			var stackData interface{} = nil
			var pipeData interface{} = nil
			if containerData, ok := stackDataMap[id]; ok {
				stackData = containerData.StackData
				pipeData = containerData.OutPipeData
			}

			obs, err = cb(stackData, pipeData, obs)
			if err != nil {
				return err
			}
		}
		return nil
	}
	for i := 0; i < len(self.Stacks); i++ {
		stack := self.Stacks[i].GetStackDefinition().GetOutbound()
		if stack != nil {
			//var err error
			//var boundResult boundResult
			boundResultInstance, err := stack.GetBoundResult()
			if boundResultInstance == nil {
				return nil, goerrors.InvalidParam
			}

			var stackBoundDefinition IStackBoundDefinition
			stackBoundDefinition, err = boundResultInstance()
			if err != nil {
				return nil, err
			}
			err = handleStack(self.Stacks[i].GetStackDefinition().GetId(), stackBoundDefinition)
			if err != nil {
				return nil, err
			}
		}
	}

	return obs, nil
}

func (self *TwoWayPipeDefinition) BuildInBoundPipeStates() ([]*PipeState, error) {
	var pipeStarts []*PipeState

	for _, currentStack := range self.Stacks {
		if currentStack == nil {
			continue
		}
		stack := currentStack.GetStackDefinition().GetInbound()
		if stack == nil {
			continue
		}
		boundResult, err := stack.GetBoundResult()
		if err != nil {
			return nil, err
		}
		if boundResult == nil {
			continue
		}
		var stackBoundDefinition IStackBoundDefinition
		stackBoundDefinition, err = boundResult()
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

func (self *TwoWayPipeDefinition) buildInBoundPipesObservables(stackDataMap map[string]*StackDataContainer, inbound chan rxgo.Item, opts ...rxgo.Option) (rxgo.Observable, error) {
	obs := rxgo.FromChannel(inbound, opts...)

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
	for i := len(self.Stacks) - 1; i >= 0; i-- {
		stack := self.Stacks[i].GetStackDefinition().GetInbound()
		if stack != nil {
			var err error
			var boundResultInstance BoundResult
			boundResultInstance, err = stack.GetBoundResult()
			if boundResultInstance == nil {
				return nil, goerrors.InvalidParam
			}

			var stackBoundDefinition IStackBoundDefinition
			stackBoundDefinition, err = boundResultInstance()
			if err != nil {
				return nil, err
			}
			err = handleStack(self.Stacks[i].GetStackDefinition().GetId(), stackBoundDefinition)
			if err != nil {
				return nil, err
			}
		}
	}
	return obs, nil
}

func NewTwoWayPipeDefinition() *TwoWayPipeDefinition {
	return &TwoWayPipeDefinition{}
}
