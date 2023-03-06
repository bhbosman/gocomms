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
					LifeCycle      fx.Lifecycle
					PipeDefinition common.IPipeDefinition `name:"Inbound"`
				},
			) ([]*common.PipeState, error) {
				return params.PipeDefinition.BuildPipeStates()
			},
		},
	)
}
