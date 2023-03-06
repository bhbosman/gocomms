package internal

import (
	"github.com/bhbosman/gocomms/common"
	"go.uber.org/fx"
)

func ProvideExtractPipeOutStates(name string) fx.Option {
	return fx.Provide(
		fx.Annotated{
			Name: name,
			Target: func(
				params struct {
					fx.In
					LifeCycle              fx.Lifecycle
					OutboundPipeDefinition common.IPipeDefinition `name:"Outbound"`
				},
			) ([]*common.PipeState, error) {
				return params.OutboundPipeDefinition.BuildOutBoundPipeStates()
			},
		},
	)
}
