package common

import (
	"context"
	"github.com/bhbosman/gocommon"
	"net"
)

type IPipeCreateData interface{}

type PipeDefinition func(
	stackData IStackCreateData,
	pipeData IPipeCreateData,
	obs gocommon.IObservable,
) (gocommon.IObservable, error)

type PipeCreate func(stackData IStackCreateData, ctx context.Context) (interface{}, error)
type PipeDestroy func(stackData IStackCreateData, pipeData IPipeCreateData) error
type PipeStart func(stackData IStackCreateData, pipeData IPipeCreateData, ctx context.Context) error
type PipeEnd func(stackData IStackCreateData, pipeData IPipeCreateData) error
type PipeState struct {
	ID      string
	Create  PipeCreate
	Destroy PipeDestroy
	Start   PipeStart
	End     PipeEnd
}

type IInputStreamForStack net.Conn

type StackEndStateParams struct {
}

type StackState struct {
	Id          string
	HijackStack bool
	Create      StackCreate
	Destroy     StackDestroy
	Start       StackStartState
	Stop        StackStopState
}

func (self *StackState) GetId() string {
	return self.Id
}

func (self *StackState) GetHijackStack() bool {
	return self.HijackStack
}

func (self *StackState) OnCreate() StackCreate {
	return self.Create
}

func (self *StackState) OnDestroy() StackDestroy {
	return self.Destroy
}

func (self *StackState) OnStart() StackStartState {
	return self.Start
}

func (self *StackState) OnStop() StackStopState {
	return self.Stop
}

type IStackState interface {
	GetId() string
	GetHijackStack() bool
	OnCreate() StackCreate
	OnDestroy() StackDestroy
	OnStart() StackStartState
	OnStop() StackStopState
}

func NewStackState(
	id string,
	hijackStack bool,
	create StackCreate,
	destroy StackDestroy,
	start StackStartState,
	stop StackStopState,
) IStackState {
	return &StackState{
		Id:          id,
		HijackStack: hijackStack,
		Create:      create,
		Destroy:     destroy,
		Start:       start,
		Stop:        stop,
	}
}

type StackDataContainer struct {
	StackData   IStackCreateData
	InPipeData  interface{}
	OutPipeData interface{}
}
