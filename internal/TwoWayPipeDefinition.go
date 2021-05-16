package internal

import (
	"context"
	"fmt"
	"github.com/bhbosman/gocomms/connectionManager"
	"github.com/bhbosman/goerrors"
	"github.com/google/uuid"
	"github.com/reactivex/rxgo/v2"
	"strings"
)

type StackDefinitionIndex struct {
	Idx             int
	StackDefinition IStackDefinition
}

type TwoWayPipeDefinition struct {
	Stacks []*StackDefinitionIndex
}

func (self *TwoWayPipeDefinition) AddStackDefinition(stack IStackDefinition) {
	index := len(self.Stacks) * 1024
	self.Stacks = append(
		self.Stacks,
		&StackDefinitionIndex{
			Idx:             index,
			StackDefinition: stack,
		})
}
func (self *TwoWayPipeDefinition) AddStackDefinitionFunc(fn func() (IStackDefinition, error)) {
	if fn == nil {
		return
	}
	definition, err := fn()
	if err != nil {
		return
	}
	self.AddStackDefinition(definition)
}

func (self TwoWayPipeDefinition) BuildStackState() ([]*StackState, error) {
	var allStackState []*StackState
	for _, item := range self.Stacks {
		stackState := item.StackDefinition.GetStackState()
		if stackState == nil {
			continue
		}
		b := true
		b = b && (0 != strings.Compare(uuid.Nil.String(), stackState.Id.String()))
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

func (self TwoWayPipeDefinition) BuildIncomingObs(
	stackDataMap map[uuid.UUID]interface{},
	connectionId string,
	connectionManager connectionManager.IConnectionManager,
	cancelCtx context.Context,
	stackCancelFunc CancelFunc) (*IncomingObs, error) {
	inBoundChannel := NewChannelManager("inboundChannelManager", connectionId)
	var obsIn rxgo.Observable
	var err error
	obsIn, err = self.buildInBoundPipesObservables(
		stackDataMap,
		connectionId,
		connectionManager,
		cancelCtx,
		stackCancelFunc,
		inBoundChannel.Items,
		rxgo.WithContext(cancelCtx))
	if err != nil {
		return nil, err
	}
	return &IncomingObs{
		InboundChannelManager: inBoundChannel,
		InboundObservable:     obsIn,
		CancelCtx:             cancelCtx,
	}, nil
}

func (self TwoWayPipeDefinition) BuildOutgoingObs(
	stackDataMap map[uuid.UUID]interface{},
	connectionId string,
	connectionManager connectionManager.IConnectionManager,
	cancelCtx context.Context,
	stackCancelFunc CancelFunc) (*OutgoingObs, error) {
	outBoundChannel := NewChannelManager("outboundChannelManager", connectionId)
	var err error
	var obsOut rxgo.Observable
	obsOut, err = self.buildOutBoundObservables(
		stackDataMap,
		connectionId,
		connectionManager,
		cancelCtx,
		stackCancelFunc,
		outBoundChannel,
		rxgo.WithContext(cancelCtx))
	if err != nil {
		return nil, err
	}

	return &OutgoingObs{
		OutboundChannelManager: outBoundChannel,
		OutboundObservable:     obsOut,
		CancelCtx:              cancelCtx,
	}, nil
}

func (self *TwoWayPipeDefinition) BuildOutBoundPipeStates() ([]*PipeState, error) {
	var pipeStarts []*PipeState

	for _, currentStack := range self.Stacks {
		if currentStack == nil {
			continue
		}
		stack := currentStack.StackDefinition.GetOutbound()
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
		stackBoundDefinition, err = boundResult(NewInOutBoundParams(currentStack.Idx))
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
		b = b && (0 != strings.Compare(uuid.Nil.String(), pipeState.ID.String()))
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

func (self *TwoWayPipeDefinition) buildOutBoundObservables(
	stackDataMap map[uuid.UUID]interface{},
	connectionId string,
	connectionManager connectionManager.IConnectionManager,
	cancelContext context.Context,
	stackCancelFunc CancelFunc,
	outbound *ChannelManager,
	opts ...rxgo.Option) (rxgo.Observable, error) {
	obs := rxgo.FromChannel(outbound.Items, opts...)

	handleStack := func(id uuid.UUID, currentStack IStackBoundDefinition) error {
		cb := currentStack.GetPipeDefinition()
		if cb != nil {
			var err error
			stackData, _ := stackDataMap[id]
			_, obs, err = cb(stackData, nil, NewPipeDefinitionParams(connectionId, connectionManager, cancelContext, stackCancelFunc, obs))
			if err != nil {
				return err
			}
		}
		return nil
	}
	for i := 0; i < len(self.Stacks); i++ {
		stack := self.Stacks[i].StackDefinition.GetOutbound()
		if stack != nil {
			var err error
			var boundResult BoundResult
			boundResult, err = stack.GetBoundResult()
			if boundResult == nil {
				return nil, goerrors.InvalidParam
			}

			var stackBoundDefinition IStackBoundDefinition
			stackBoundDefinition, err = boundResult(
				NewInOutBoundParams(self.Stacks[i].Idx))
			if err != nil {
				return nil, err
			}
			err = handleStack(self.Stacks[i].StackDefinition.GetId(), stackBoundDefinition)
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
		stack := currentStack.StackDefinition.GetInbound()
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
		stackBoundDefinition, err = boundResult(NewInOutBoundParams(currentStack.Idx))
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
		b = b && (0 != strings.Compare(uuid.Nil.String(), pipeState.ID.String()))
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

func (self *TwoWayPipeDefinition) buildInBoundPipesObservables(
	stackDataMap map[uuid.UUID]interface{},
	connectionId string,
	connectionManager connectionManager.IConnectionManager,
	cancelContext context.Context,
	stackCancelFunc CancelFunc,
	inbound chan rxgo.Item,
	opts ...rxgo.Option) (rxgo.Observable, error) {
	obs := rxgo.FromChannel(inbound, opts...)

	handleStack := func(id uuid.UUID, currentStack IStackBoundDefinition) error {
		cb := currentStack.GetPipeDefinition()
		if cb != nil {
			stackData, _ := stackDataMap[id]

			var err error
			_, obs, err = cb(stackData, nil, NewPipeDefinitionParams(connectionId, connectionManager, cancelContext, stackCancelFunc, obs))
			if err != nil {
				return err
			}
		}
		return nil
	}
	for i := len(self.Stacks) - 1; i >= 0; i-- {
		stack := self.Stacks[i].StackDefinition.GetInbound()
		if stack != nil {
			var err error
			var boundResult BoundResult
			boundResult, err = stack.GetBoundResult()
			if boundResult == nil {
				return nil, goerrors.InvalidParam
			}

			var stackBoundDefinition IStackBoundDefinition
			stackBoundDefinition, err = boundResult(
				NewInOutBoundParams(self.Stacks[i].Idx))
			if err != nil {
				return nil, err
			}
			err = handleStack(self.Stacks[i].StackDefinition.GetId(), stackBoundDefinition)
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
