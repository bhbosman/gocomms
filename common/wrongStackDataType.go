package common

import (
	"fmt"
	"github.com/bhbosman/gocommon/model"
	"reflect"
)

type WrongStackDataType struct {
	StackName      string
	ConnectionType model.ConnectionType
	WantedType     reflect.Type
	ReceivedType   reflect.Type
}

func NewWrongStackDataType(
	stackName string,
	connectionType model.ConnectionType,
	wantedType reflect.Type,
	receivedType reflect.Type,
) *WrongStackDataType {
	return &WrongStackDataType{
		StackName:      stackName,
		ConnectionType: connectionType,
		WantedType:     wantedType,
		ReceivedType:   receivedType,
	}
}

func (w WrongStackDataType) Error() string {
	return fmt.Sprintf("Wrong data type")
}
