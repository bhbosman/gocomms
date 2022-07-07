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
					LifeCycle            fx.Lifecycle
					TwoWayPipeDefinition *common.TwoWayPipeDefinition
				},
			) ([]*common.PipeState, error) {
				return params.TwoWayPipeDefinition.BuildOutBoundPipeStates()
			},
		},
	)
}
