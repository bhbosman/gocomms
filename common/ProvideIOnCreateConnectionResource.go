package common

import (
	"github.com/bhbosman/goCommsDefinitions"
)

func ProvideIOnCreateConnectionResource() (goCommsDefinitions.IOnCreateConnection, error) {
	return &OnCreateConnectionDefault{}, nil
}
