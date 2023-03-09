package common

import (
	"github.com/bhbosman/gocommon/model"
	"github.com/reactivex/rxgo/v2"
)

type IStackCreateData interface{}
type StackCreate func() (IStackCreateData, error)

type StackDestroy func(
	connectionType model.ConnectionType,
	stackData IStackCreateData) error

type StackStartState func(
	conn IInputStreamForStack,
	stackData IStackCreateData,
	ToReactorFunc rxgo.NextFunc,
) (IInputStreamForStack, error)

type StackStopState func(
	stackData interface{},
) error

type stackState struct {
	Id          string
	HijackStack bool
	Create      StackCreate
	Destroy     StackDestroy
	Start       StackStartState
	Stop        StackStopState
}

func (self *stackState) GetId() string {
	return self.Id
}

func (self *stackState) GetHijackStack() bool {
	return self.HijackStack
}

func (self *stackState) OnCreate() StackCreate {
	return self.Create
}

func (self *stackState) OnDestroy() StackDestroy {
	return self.Destroy
}

func (self *stackState) OnStart() StackStartState {
	return self.Start
}

func (self *stackState) OnStop() StackStopState {
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
	return &stackState{
		Id:          id,
		HijackStack: hijackStack,
		Create:      create,
		Destroy:     destroy,
		Start:       start,
		Stop:        stop,
	}
}
