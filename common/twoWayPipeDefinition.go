package common

import (
	"context"
	"fmt"
	"github.com/bhbosman/gocommon"
	"github.com/reactivex/rxgo/v2"
	"strings"
)

type outboundPipeDefinition struct {
}
type inboundPipeDefinition struct {
}

type twoWayPipeDefinition struct {
	outboundPipeDefinition outboundPipeDefinition
	inboundPipeDefinition  inboundPipeDefinition
	Stacks                 []IStackDefinition
}

func (self *twoWayPipeDefinition) BuildStackState() ([]*StackState, error) {
	var allStackState []*StackState
	for _, item := range self.Stacks {
		stackState := item.StackState()
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

func (self *twoWayPipeDefinition) BuildIncomingObs(
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

func (self *twoWayPipeDefinition) BuildOutgoingObs(
	outBoundChannel chan rxgo.Item,
	stackDataMap map[string]*StackDataContainer,
	cancelCtx context.Context,
) (*OutgoingObs, error) {
	var err error
	var obsOut gocommon.IObservable
	obsOut, err = self.buildOutBoundObservables(stackDataMap, outBoundChannel, rxgo.WithContext(cancelCtx))
	if err != nil {
		return nil, err
	}

	return NewOutgoingObs(obsOut, cancelCtx), nil
}

func (self *twoWayPipeDefinition) BuildOutBoundPipeStates() ([]*PipeState, error) {
	var pipeStarts []*PipeState

	for _, currentStack := range self.Stacks {
		if currentStack == nil {
			continue
		}
		stack := currentStack.Outbound()
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

func (self *twoWayPipeDefinition) buildOutBoundObservables(
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
		stack := self.Stacks[i].Outbound()
		if stack != nil {
			stackBoundDefinition, err := stack()
			if err != nil {
				return nil, err
			}
			err = handleStack(self.Stacks[i].Name(), stackBoundDefinition)
			if err != nil {
				return nil, err
			}
		}
	}

	return obs, nil
}

func (self *twoWayPipeDefinition) BuildInBoundPipeStates() ([]*PipeState, error) {
	var pipeStarts []*PipeState

	for _, currentStack := range self.Stacks {
		if currentStack == nil {
			continue
		}
		stack := currentStack.Inbound()
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

func (self *twoWayPipeDefinition) buildInBoundPipesObservables(
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
	for i := len(self.Stacks) - 1; i >= 0; i-- {
		stack := self.Stacks[i].Inbound()
		if stack != nil {
			stackBoundDefinition, err := stack()
			if err != nil {
				return nil, err
			}
			err = handleStack(self.Stacks[i].Name(), stackBoundDefinition)
			if err != nil {
				return nil, err
			}
		}
	}
	return obs, nil
}

type IInboundPipeDefinition interface {
	BuildOutBoundPipeStates() ([]*PipeState, error)
	BuildIncomingObs(
		inBoundChannel chan rxgo.Item,
		stackDataMap map[string]*StackDataContainer,
		cancelCtx context.Context,
	) (*IncomingObs, error)
}

type IOutboundPipeDefinition interface {
	BuildInBoundPipeStates() ([]*PipeState, error)
	BuildOutgoingObs(
		outBoundChannel chan rxgo.Item,
		stackDataMap map[string]*StackDataContainer,
		cancelCtx context.Context,
	) (*OutgoingObs, error)
}

type ITwoWayPipeDefinition interface {
	IOutboundPipeDefinition
	IInboundPipeDefinition
	BuildStackState() ([]*StackState, error)
}

func NewTwoWayPipeDefinition(stacks []IStackDefinition) ITwoWayPipeDefinition {
	return &twoWayPipeDefinition{
		outboundPipeDefinition: outboundPipeDefinition{},
		inboundPipeDefinition:  inboundPipeDefinition{},
		Stacks:                 stacks,
	}
}
