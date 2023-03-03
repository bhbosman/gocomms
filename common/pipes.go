package common

import (
	"context"
	"github.com/bhbosman/gocommon/model"
	"github.com/reactivex/rxgo/v2"
	"net"
)

type IPipeCreateData interface{}

//type IObservable interface {
//	rxgo.Observable
//}

type PipeDefinition func(stackData IStackCreateData, pipeData IPipeCreateData, obs rxgo.Observable) (rxgo.Observable, error)

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

type IStackCreateData interface{}
type StackCreate func() (IStackCreateData, error)

type StackDestroy func(
	connectionType model.ConnectionType,
	stackData IStackCreateData) error

type IInputStreamForStack net.Conn

type StackStartState func(
	conn IInputStreamForStack,
	stackData IStackCreateData,
	ToReactorFunc rxgo.NextFunc,
) (IInputStreamForStack, error)
type StackStopState func(stackData interface{}, endParams StackEndStateParams) error

type StackEndStateParams struct {
}

func NewStackEndStateParams() StackEndStateParams {
	return StackEndStateParams{}
}

type StackState struct {
	Id          string
	HijackStack bool
	Create      StackCreate
	Destroy     StackDestroy
	Start       StackStartState
	Stop        StackStopState
}

type StackDataContainer struct {
	StackData   IStackCreateData
	InPipeData  interface{}
	OutPipeData interface{}
}
