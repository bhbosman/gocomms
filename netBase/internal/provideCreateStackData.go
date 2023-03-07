package internal

import (
	"github.com/bhbosman/gocommon/model"
	"github.com/bhbosman/gocomms/common"
	"github.com/bhbosman/gocomms/intf"
	"go.uber.org/fx"
	"golang.org/x/net/context"
	"net"
	"net/url"
)

func ProvideCreateStackData() fx.Option {
	return fx.Provide(
		func(
			params struct {
				fx.In
				LifeCycle         fx.Lifecycle
				StackState        []common.IStackState
				PipeInStates      []*common.PipeState `name:"PipeInStates"`
				PipeOutStates     []*common.PipeState `name:"PipeOutStates"`
				Conn              net.Conn            `name:"PrimaryConnection"`
				Url               *url.URL
				Ctx               context.Context
				CtxCancelFunc     context.CancelFunc
				CancelFunc        model.ConnectionCancelFunc
				Ev                intf.IConnectionReactorFactoryExtractValues `optional:"true"`
				ConnectionReactor intf.IConnectionReactor
				ConnectionType    model.ConnectionType
			},
		) (map[string]*common.StackDataContainer, error) {

			// to make sure all usages will have some value
			if params.Ev == nil {
				params.Ev = intf.NewDefaultConnectionReactorFactoryExtractValues()
			}

			result := make(map[string]*common.StackDataContainer)
			var err error
			for _, iteration := range params.StackState {
				localItem := iteration
				var stackData interface{}
				stackData, err = localItem.OnCreate()()
				if err != nil {
					return nil, err
				}
				if stackData == nil {
					stackData = common.NewNoCloser()
				}
				if container, ok := result[localItem.GetId()]; ok {
					container.StackData = stackData
				} else {
					result[localItem.GetId()] = &common.StackDataContainer{
						StackData: stackData,
					}
				}
				params.LifeCycle.Append(
					fx.Hook{
						OnStop: func(ctx context.Context) error {
							if localItem.OnDestroy() != nil {
								if container, ok := result[localItem.GetId()]; ok {
									return localItem.OnDestroy()(params.ConnectionType, container.StackData)
								}
							}
							return nil
						},
					})
			}
			for _, iteration := range params.PipeInStates {
				localItem := iteration
				var stackData interface{} = nil
				if container, ok := result[localItem.ID]; ok {
					stackData = container.StackData
				}

				var pipeData interface{}
				pipeData, err = localItem.Create(stackData, params.Ctx)
				if err != nil {
					return nil, err
				}
				if pipeData == nil {
					pipeData = common.NewNoCloser()
				}
				if container, ok := result[localItem.ID]; ok {
					container.InPipeData = pipeData
				} else {
					result[localItem.ID] = &common.StackDataContainer{
						InPipeData: pipeData,
					}
				}
				params.LifeCycle.Append(
					fx.Hook{
						OnStop: func(ctx context.Context) error {
							if localItem.Destroy != nil {
								if container, ok := result[localItem.ID]; ok {
									return localItem.Destroy(container.StackData, container.InPipeData)
								}
							}
							return nil
						},
					})
			}
			for _, iteration := range params.PipeOutStates {
				localItem := iteration
				var stackData interface{} = nil
				if container, ok := result[localItem.ID]; ok {
					stackData = container.StackData
				}

				var pipeData interface{}
				pipeData, err = localItem.Create(stackData, params.Ctx)
				if err != nil {
					return nil, err
				}
				if pipeData == nil {
					pipeData = common.NewNoCloser()
				}
				if container, ok := result[localItem.ID]; ok {
					container.OutPipeData = pipeData
				} else {
					result[localItem.ID] = &common.StackDataContainer{
						OutPipeData: pipeData,
					}
				}
				params.LifeCycle.Append(
					fx.Hook{
						OnStop: func(ctx context.Context) error {
							if localItem.Destroy != nil {
								if container, ok := result[localItem.ID]; ok {
									return localItem.Destroy(container.StackData, container.OutPipeData)
								}
							}
							return nil
						},
					})
			}

			return result, nil
		},
	)
}
