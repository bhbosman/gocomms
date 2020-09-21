package internal

import (
	"context"
	"github.com/bhbosman/gocomms/connectionManager"
	"github.com/bhbosman/gologging"
	"github.com/reactivex/rxgo/v2"
)

type StackDefinitionIndex struct {
	Idx             int
	StackDefinition *StackDefinition
}

type TwoWayPipeDefinition struct {
	Stacks []*StackDefinitionIndex
}

func (self *TwoWayPipeDefinition) Add(name string, inbound, outbound PipeDefinition) {
	self.AddStackDefinition(
		&StackDefinition{
			Inbound: func(index int, ctx context.Context) BoundDefinition {
				return BoundDefinition{
					PipeDefinition: inbound,
				}
			},
			Outbound: func(index int, ctx context.Context) BoundDefinition {
				return BoundDefinition{
					PipeDefinition: outbound,
				}
			},
			Name: name,
		})
}

func (self *TwoWayPipeDefinition) AddStackDefinition(stack *StackDefinition) {
	index := len(self.Stacks) * 1024
	self.Stacks = append(
		self.Stacks,
		&StackDefinitionIndex{
			Idx:             index,
			StackDefinition: stack,
		})
}
func (self *TwoWayPipeDefinition) AddStackDefinitionFunc(fn func() (*StackDefinition, error)) {
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
	logger *gologging.SubSystemLogger,
	cancelCtx context.Context,
	stackCancelFunc CancelFunc) (*TwoWayPipe, error) {
	createChannel := func() chan rxgo.Item {
		return make(chan rxgo.Item)
	}
	inBoundChannel := createChannel()
	outBoundChannel := createChannel()

	var allIndividualState, pipeState []PipeState
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
		allStackState = append(allStackState, stack.StackDefinition.StackState)
	}

	return NewTwoWayPipe(
		logger,
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
	opts ...rxgo.Option) (rxgo.Observable, []PipeState, error) {
	obs := rxgo.FromChannel(outbound, opts...)
	var pipeStarts []PipeState
	handleStack := func(currentStack BoundDefinition) error {
		if currentStack.PipeDefinition != nil {
			var err error
			obs, err = currentStack.PipeDefinition(NewPipeDefinitionParams(connectionId, connectionManager, cancelContext, stackCancelFunc, obs))
			if err != nil {
				return err
			}
		}
		pipeStarts = append(pipeStarts, currentStack.PipeState)

		return nil
	}
	for i := 0; i < len(self.Stacks); i++ {
		stack := self.Stacks[i].StackDefinition.Outbound
		if stack != nil {
			err := handleStack(stack(self.Stacks[i].Idx, cancelContext))
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
	opts ...rxgo.Option) (rxgo.Observable, []PipeState, error) {
	obs := rxgo.FromChannel(inbound, opts...)
	var pipeStarts []PipeState

	handleStack := func(currentStack BoundDefinition) error {
		if currentStack.PipeDefinition != nil {
			var err error
			obs, err = currentStack.PipeDefinition(NewPipeDefinitionParams(connectionId, connectionManager, cancelContext, stackCancelFunc, obs))
			if err != nil {
				return err
			}
		}
		pipeStarts = append(pipeStarts, currentStack.PipeState)
		return nil
	}
	for i := len(self.Stacks) - 1; i >= 0; i-- {
		stack := self.Stacks[i].StackDefinition.Inbound
		if stack != nil {
			err := handleStack(stack(self.Stacks[i].Idx, cancelContext))
			if err != nil {
				return nil, nil, err
			}
		}
	}
	return obs, pipeStarts, nil
}

func NewTwoWayPipeDefinition(Stacks []*StackDefinitionIndex) *TwoWayPipeDefinition {
	return &TwoWayPipeDefinition{
		Stacks: Stacks,
	}
}
