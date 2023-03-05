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
)

func ProvideCreateStackDefinition() fx.Option {
	return fx.Provide(
		fx.Annotated{
			Target: func(
				params struct {
					fx.In
					CancelFunc           context.CancelFunc
					ConnectionCancelFunc model.ConnectionCancelFunc
					Logger               *zap.Logger
					StackName            string                                 `name:"StackName"`
					TransportFactories   []*goCommsDefinitions.TransportFactory `group:"TransportFactory"`
					StackFactories       []common.IStackDefinition              `group:"StackDefinition"`
				},
			) (common.ITwoWayPipeDefinition, common.IInboundPipeDefinition, common.IOutboundPipeDefinition, error) {
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
						return nil, nil, nil, errList
					}

					var stacks []common.IStackDefinition
					for _, stackName := range factory.StackNames {
						if item, ok := dict[stackName]; ok {
							stacks = append(stacks, item)
						}
					}
					inboundPipeDefinition := common.NewInboundPipeDefinition(stacks)
					outboundPipeDefinition := common.NewOutboundPipeDefinition(stacks)
					twoWayPipeDefinition, err := common.NewTwoWayPipeDefinition(
						stacks,
						outboundPipeDefinition,
						inboundPipeDefinition)
					if err != nil {
						return nil, nil, nil, err
					}

					return twoWayPipeDefinition, inboundPipeDefinition, outboundPipeDefinition, nil
				}
				return nil, nil, nil, fmt.Errorf("connectionstack factory definition \"%v\", not found", params.StackName)
			},
		},
	)
}
