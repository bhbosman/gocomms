package internal

import (
	"context"
	"fmt"
	"github.com/bhbosman/goCommsDefinitions"
	"github.com/bhbosman/gocommon/model"
	"github.com/bhbosman/gocomms/common"
	"go.uber.org/fx"
	"go.uber.org/multierr"
	"go.uber.org/zap"
	"strings"
)

type CreateStackDefinition struct {
	fx.Out
	StackState             []common.IStackState
	InboundPipeDefinition  common.IInboundPipeDefinition
	OutboundPipeDefinition common.IPipeDefinition `name:"Outbound"`
}

func ProvideCreateStackDefinition() fx.Option {
	return fx.Provide(

		func(
			params struct {
				fx.In
				CancelFunc           context.CancelFunc
				ConnectionCancelFunc model.ConnectionCancelFunc
				Logger               *zap.Logger
				StackName            string                                 `name:"StackName"`
				TransportFactories   []*goCommsDefinitions.TransportFactory `group:"TransportFactory"`
				StackFactories       []common.IStackDefinition              `group:"StackDefinition"`
			},
		) (CreateStackDefinition, error) {
			params.Logger.Info("createStackDefinition...")
			var factory *goCommsDefinitions.TransportFactory = nil
			for _, item := range params.TransportFactories {
				if item.Name == params.StackName {
					factory = item
					break
				}
			}
			if factory != nil {
				dict := make(map[string]common.IStackDefinition)
				for _, stackFactory := range params.StackFactories {
					dict[stackFactory.Name()] = stackFactory
				}
				var errList error = nil
				for _, stackName := range factory.StackNames {
					if _, ok := dict[stackName]; !ok {
						s := fmt.Sprintf("stack definition %v, not found", stackName)
						errList = multierr.Append(errList, fmt.Errorf(s))
					}
				}
				if errList != nil {
					params.CancelFunc()
					params.ConnectionCancelFunc("On stack creation", false, errList)
					return CreateStackDefinition{}, errList
				}

				inboundPipeDefinition := common.NewInboundPipeDefinition(
					func() []common.IBoundData {
						var result []common.IBoundData
						for i := len(factory.StackNames) - 1; i >= 0; i-- {
							stackName := factory.StackNames[i]
							if item, ok := dict[stackName]; ok {
								result = append(result,
									common.NewBoundData(
										item.Name(),
										item.Inbound(),
									),
								)
							}
						}
						return result
					}())

				outboundPipeDefinition := common.NewPipeDefinition(
					func() []common.IBoundData {
						var result []common.IBoundData
						for _, stackName := range factory.StackNames {
							if item, ok := dict[stackName]; ok {
								result = append(result,
									common.NewBoundData(
										item.Name(),
										item.Outbound(),
									),
								)
							}
						}
						return result
					}())

				stackState, err := func() ([]common.IStackState, error) {
					var allStackState []common.IStackState
					var stacks []common.IStackDefinition
					for _, stackName := range factory.StackNames {
						if item, ok := dict[stackName]; ok {
							stacks = append(stacks, item)
						}
					}
					for _, item := range stacks {
						stackState := item.StackState()
						if stackState == nil {
							continue
						}
						b := true
						b = b && (0 != strings.Compare("", stackState.GetId()))
						b = b && (stackState.OnCreate() != nil)
						b = b && (stackState.OnDestroy() != nil)
						b = b && (stackState.OnStart() != nil)
						b = b && (stackState.OnStop() != nil)
						if !b {
							return nil, fmt.Errorf("stackstate must be complete in full")
						}
						allStackState = append(allStackState, stackState)
					}
					return allStackState, nil
				}()
				if err != nil {
					return CreateStackDefinition{}, err
				}

				return CreateStackDefinition{
					StackState:             stackState,
					InboundPipeDefinition:  inboundPipeDefinition,
					OutboundPipeDefinition: outboundPipeDefinition,
				}, nil
			}
			return CreateStackDefinition{}, fmt.Errorf("connectionstack factory definition \"%v\", not found", params.StackName)
		},
	)
}
