package internal

import (
	"github.com/bhbosman/gocomms/common"
	"go.uber.org/fx"
)

func ProvideExtractPipeInStates(name string) fx.Option {
	return fx.Provide(
		fx.Annotated{
			Name: name,
			Target: func(
				params struct {
					fx.In
					LifeCycle             fx.Lifecycle
					InboundPipeDefinition common.IInboundPipeDefinition
				},
			) ([]*common.PipeState, error) {
				return params.InboundPipeDefinition.BuildInBoundPipeStates()
			},
		},
	)
}
