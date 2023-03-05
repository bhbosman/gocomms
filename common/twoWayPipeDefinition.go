package common

import (
	"context"
	"fmt"
	"github.com/reactivex/rxgo/v2"
	"strings"
)

type twoWayPipeDefinition struct {
	outboundPipeDefinition outboundPipeDefinition
	inboundPipeDefinition  inboundPipeDefinition
	stacks                 []IStackDefinition
}

func (self *twoWayPipeDefinition) BuildStackState() ([]*StackState, error) {
	var allStackState []*StackState
	for _, item := range self.stacks {
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
	return self.inboundPipeDefinition.BuildIncomingObs(inBoundChannel, stackDataMap, cancelCtx)
}

func (self *twoWayPipeDefinition) BuildOutgoingObs(
	outBoundChannel chan rxgo.Item,
	stackDataMap map[string]*StackDataContainer,
	cancelCtx context.Context,
) (*OutgoingObs, error) {
	return self.outboundPipeDefinition.BuildOutgoingObs(outBoundChannel, stackDataMap, cancelCtx)
}

func (self *twoWayPipeDefinition) BuildOutBoundPipeStates() ([]*PipeState, error) {
	var pipeStarts []*PipeState

	for _, currentStack := range self.stacks {
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

func (self *twoWayPipeDefinition) BuildInBoundPipeStates() ([]*PipeState, error) {
	var pipeStarts []*PipeState

	for _, currentStack := range self.stacks {
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
		outboundPipeDefinition: outboundPipeDefinition{
			stacks: stacks,
		},
		inboundPipeDefinition: inboundPipeDefinition{
			stacks: stacks,
		},
		stacks: stacks,
	}
}
