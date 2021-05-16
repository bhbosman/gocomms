package internal

import (
	"fmt"
	"reflect"
)

type NoCloser struct {
}

func NewNoCloser() *NoCloser {
	return &NoCloser{}
}

func (n NoCloser) Close() error {
	return nil
}

type IStackBoundDefinition interface {
	GetPipeDefinition() PipeDefinition
	GetPipeState() *PipeState
}

type StackBoundDefinition struct {
	PipeDefinition PipeDefinition
	PipeState      *PipeState
}

func (self *StackBoundDefinition) GetPipeDefinition() PipeDefinition {
	return self.PipeDefinition
}

func (self *StackBoundDefinition) GetPipeState() *PipeState {
	return self.PipeState
}

func NewBoundDefinition(pipeDefinition PipeDefinition, pipeState *PipeState) IStackBoundDefinition {
	return &StackBoundDefinition{
		PipeDefinition: pipeDefinition,
		PipeState:      pipeState,
	}
}

type WrongStackDataType struct {
	StackName      string
	ConnectionType ConnectionType
	WantedType     reflect.Type
	ReceivedType   reflect.Type
}

func NewWrongStackDataType(stackName string, connectionType ConnectionType, wantedType reflect.Type, receivedType reflect.Type) *WrongStackDataType {
	return &WrongStackDataType{StackName: stackName, ConnectionType: connectionType, WantedType: wantedType, ReceivedType: receivedType}
}

func (w WrongStackDataType) Error() string {
	return fmt.Sprintf("Wrong data type")
}
