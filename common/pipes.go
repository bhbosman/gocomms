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
	ID        string
	OnCreate  PipeCreate
	OnDestroy PipeDestroy
	OnStart   PipeStart
	OnEnd     PipeEnd
}

type IInputStreamForStack net.Conn

type StackEndStateParams struct {
}

type StackDataContainer struct {
	StackData   IStackCreateData
	InPipeData  interface{}
	OutPipeData interface{}
}
