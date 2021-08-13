package internal

import (
	"context"
	"github.com/bhbosman/gocomms/internal"
	"github.com/google/uuid"
	"go.uber.org/fx"
)

func CreatePipeStateMap(
	Ctx context.Context,
	PipeInData []*internal.PipeState,
	StackData map[uuid.UUID]interface{},
	LifeCycle fx.Lifecycle) (map[uuid.UUID]interface{}, error) {
	result := make(map[uuid.UUID]interface{})
	for _, iteration := range PipeInData {
		if iteration == nil {
			continue
		}
		localItem := iteration
		stackData, _ := StackData[localItem.ID]
		pipeData, err := localItem.Create(stackData, Ctx)
		if err != nil {
			return nil, err
		}
		if pipeData == nil {
			pipeData = internal.NewNoCloser()
		}
		result[localItem.ID] = pipeData
		LifeCycle.Append(fx.Hook{
			OnStart: nil,
			OnStop: func(ctx context.Context) error {
				if localItem.Destroy != nil {
					stackData, _ := StackData[localItem.ID]
					pipeData, _ := result[localItem.ID]
					return localItem.Destroy(stackData, pipeData)
				}
				return nil
			},
		})
	}

	return result, nil
}
