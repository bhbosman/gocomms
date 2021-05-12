package internal

import (
	"context"
	"github.com/bhbosman/gocomms/connectionManager"
	"github.com/bhbosman/goerrors"
	"github.com/google/uuid"
	"github.com/reactivex/rxgo/v2"
	"net"
)

type StackDefinitionIndex struct {
	Idx             int
	StackDefinition IStackDefinition
}

type TwoWayPipeDefinition struct {
	Stacks []*StackDefinitionIndex
}

func (self *TwoWayPipeDefinition) Add(id uuid.UUID, name string, inbound, outbound PipeDefinition) {
	self.AddStackDefinition(
		&StackDefinition{
			Inbound: NewBoundResultImpl(
				func(stackData, pipeData interface{}, params InOutBoundParams) IStackBoundDefinition {
					return NewBoundDefinition(inbound, nil)
				}),
			Outbound: NewBoundResultImpl(
				func(stackData, pipeData interface{}, params InOutBoundParams) IStackBoundDefinition {
					return NewBoundDefinition(outbound, nil)

				}),
			Name: name,
			Id:   id,
		})
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

func (self TwoWayPipeDefinition) Build(
	connectionId string,
	connectionManager connectionManager.IConnectionManager,
	conn net.Conn,
	cancelCtx context.Context,
	stackCancelFunc CancelFunc) (*TwoWayPipe, error) {
	createChannel := func() chan rxgo.Item {
		return make(chan rxgo.Item, 1024)
	}
	inBoundChannel := createChannel()
	outBoundChannel := createChannel()

	var allIndividualState, pipeState []*PipeState
	var allStackState []StackState
	var err error

	var obsOut rxgo.Observable
	obsOut, pipeState, err = self.buildOutBound(
		connectionId,
		connectionManager,
		cancelCtx,
		stackCancelFunc,
		outBoundChannel, rxgo.WithContext(cancelCtx))
	if err != nil {
		return nil, err
	}
	allIndividualState = append(allIndividualState, pipeState...)

	var obsIn rxgo.Observable
	obsIn, pipeState, err = self.buildInBound(
		connectionId,
		connectionManager,
		cancelCtx,
		stackCancelFunc,
		inBoundChannel,
		rxgo.WithContext(cancelCtx))
	if err != nil {
		return nil, err
	}
	allIndividualState = append(allIndividualState, pipeState...)

	for i := len(self.Stacks) - 1; i >= 0; i-- {
		stack := self.Stacks[i]
		allStackState = append(allStackState, stack.StackDefinition.GetStackState())
	}

	return NewTwoWayPipe(
		connectionId,
		//logger,
		inBoundChannel,
		outBoundChannel,
		obsIn,
		obsOut,
		cancelCtx,
		allIndividualState,
		allStackState), nil
}

func (self *TwoWayPipeDefinition) buildOutBound(
	connectionId string,
	connectionManager connectionManager.IConnectionManager,
	cancelContext context.Context,
	stackCancelFunc CancelFunc,
	outbound chan rxgo.Item,
	opts ...rxgo.Option) (rxgo.Observable, []*PipeState, error) {
	obs := rxgo.FromChannel(outbound, opts...)
	var pipeStarts []*PipeState
	handleStack := func(currentStack IStackBoundDefinition) error {
		cb := currentStack.GetPipeDefinition()
		if cb != nil {
			var err error
			_, obs, err = cb(NewPipeDefinitionParams(connectionId, connectionManager, cancelContext, stackCancelFunc, obs))
			if err != nil {
				return err
			}
		}
		pipeStarts = append(pipeStarts, currentStack.GetPipeState())

		return nil
	}
	for i := 0; i < len(self.Stacks); i++ {
		stack := self.Stacks[i].StackDefinition.GetOutbound()
		if stack != nil {
			cb := stack.GetBoundResult()
			if cb == nil {
				return nil, nil, goerrors.InvalidParam
			}
			err := handleStack(cb(
				nil,
				nil,
				NewInOutBoundParams(
					self.Stacks[i].Idx,
					cancelContext,
				)))
			if err != nil {
				return nil, nil, err
			}
		}
	}

	return obs, pipeStarts, nil
}
func (self *TwoWayPipeDefinition) buildInBound(
	connectionId string,
	connectionManager connectionManager.IConnectionManager,
	cancelContext context.Context,
	stackCancelFunc CancelFunc,
	inbound chan rxgo.Item,
	opts ...rxgo.Option) (rxgo.Observable, []*PipeState, error) {
	obs := rxgo.FromChannel(inbound, opts...)
	var pipeStarts []*PipeState

	handleStack := func(currentStack IStackBoundDefinition) error {
		cb := currentStack.GetPipeDefinition()
		if cb != nil {
			var err error
			_, obs, err = cb(NewPipeDefinitionParams(connectionId, connectionManager, cancelContext, stackCancelFunc, obs))
			if err != nil {
				return err
			}
		}
		pipeStarts = append(pipeStarts, currentStack.GetPipeState())
		return nil
	}
	for i := len(self.Stacks) - 1; i >= 0; i-- {
		stack := self.Stacks[i].StackDefinition.GetInbound()
		if stack != nil {
			cb := stack.GetBoundResult()
			if cb == nil {
				return nil, nil, goerrors.InvalidParam
			}
			err := handleStack(
				cb(
					nil,
					nil,
					NewInOutBoundParams(
						self.Stacks[i].Idx,
						cancelContext,
					)))
			if err != nil {
				return nil, nil, err
			}
		}
	}
	return obs, pipeStarts, nil
}

func NewTwoWayPipeDefinition() *TwoWayPipeDefinition {
	return &TwoWayPipeDefinition{}
}
