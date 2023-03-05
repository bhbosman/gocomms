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
	IOutboundPipeDefinition
	IInboundPipeDefinition

	BuildStackState() ([]*StackState, error)
}

type twoWayPipeDefinition struct {
	outboundPipeDefinition IOutboundPipeDefinition
	inboundPipeDefinition  IInboundPipeDefinition
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

func (self *twoWayPipeDefinition) BuildOutBoundPipeStates() ([]*PipeState, error) {
	return self.outboundPipeDefinition.BuildOutBoundPipeStates()
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

func (self *twoWayPipeDefinition) BuildInBoundPipeStates() ([]*PipeState, error) {
	return self.inboundPipeDefinition.BuildInBoundPipeStates()
}

func NewTwoWayPipeDefinition(
	stacks []IStackDefinition,
	outboundPipeDefinition IOutboundPipeDefinition,
	inboundPipeDefinition IInboundPipeDefinition,
) (ITwoWayPipeDefinition, error) {
	return &twoWayPipeDefinition{
		outboundPipeDefinition: outboundPipeDefinition,
		inboundPipeDefinition:  inboundPipeDefinition,
		stacks:                 stacks,
	}, nil
}
