package common

import (
	"context"
	"fmt"
	"github.com/reactivex/rxgo/v2"
	"strings"
)

type IInboundPipeDefinition interface {
	BuildInBoundPipeStates() ([]*PipeState, error)
	BuildIncomingObs(
		inBoundChannel chan rxgo.Item,
		stackDataMap map[string]*StackDataContainer,
		cancelCtx context.Context,
	) (*IncomingObs, error)
}

type IOutboundPipeDefinition interface {
	BuildOutBoundPipeStates() ([]*PipeState, error)
	BuildOutgoingObs(
		outBoundChannel chan rxgo.Item,
		stackDataMap map[string]*StackDataContainer,
		cancelCtx context.Context,
	) (*OutgoingObs, error)
}

type ITwoWayPipeDefinition interface {
	BuildStackState() ([]IStackState, error)
}

type twoWayPipeDefinition struct {
	stacks []IStackDefinition
}

func (self *twoWayPipeDefinition) BuildStackState() ([]IStackState, error) {
	var allStackState []IStackState
	for _, item := range self.stacks {
		stackState := item.StackState()
		if stackState == nil {
			continue
		}
		b := true
		b = b && (0 != strings.Compare("", stackState.GetId()))
		b = b && (stackState.OnCreate() != nil)
		b = b && (stackState.OnDestroy() != nil)
		b = b && (stackState.OnStart() != nil)
		b = b && (stackState.OnStop() != nil)
		if !b {
			return nil, fmt.Errorf("stackstate must be complete in full")
		}
		allStackState = append(allStackState, stackState)
	}
	return allStackState, nil
}

func NewTwoWayPipeDefinition(
	stacks []IStackDefinition,
) (ITwoWayPipeDefinition, error) {
	return &twoWayPipeDefinition{
		stacks: stacks,
	}, nil
}
