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
	endParams StackEndStateParams,
) error
