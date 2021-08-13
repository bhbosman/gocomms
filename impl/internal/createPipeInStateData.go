package internal

import (
	"context"
	"github.com/bhbosman/gocomms/internal"
	"github.com/google/uuid"
	"go.uber.org/fx"
)

func CreatePipeInStateData(params struct {
	fx.In
	Ctx        context.Context
	PipeInData []*internal.PipeState     `name:"PipeInStates"`
	StackData  map[uuid.UUID]interface{} `name:"StackData"`
	LifeCycle  fx.Lifecycle
}) (map[uuid.UUID]interface{}, error) {
	return CreatePipeStateMap(params.Ctx, params.PipeInData, params.StackData, params.LifeCycle)
}
